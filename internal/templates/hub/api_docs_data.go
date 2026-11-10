package hub

import "strings"

type apiParam struct {
	Name     string
	In       string // "path", "query" or "body"
	Type     string
	Required bool
	Desc     string
}

type apiResponse struct {
	Status int
	Desc   string
	Body   string // example JSON, empty if none
}

type apiEndpoint struct {
	Method    string
	Path      string
	Summary   string
	Auth      bool
	Params    []apiParam
	Request   string // example request body, empty if none
	Responses []apiResponse
	Notes     []string
}

type apiGroup struct {
	Title     string
	Endpoints []apiEndpoint
}

var errUnauthorized = apiResponse{401, "Missing, invalid or expired token.", `{ "error": "invalid or expired token" }`}

func errResp(status int, desc, message string) apiResponse {
	return apiResponse{status, desc, `{ "error": "` + message + `" }`}
}

const okBody = `{ "success": true }`

func dedent(s string) string {
	lines := strings.Split(strings.TrimPrefix(s, "\n"), "\n")
	indent := -1
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		n := len(l) - len(strings.TrimLeft(l, "\t"))
		if indent < 0 || n < indent {
			indent = n
		}
	}
	for i, l := range lines {
		if len(l) >= indent {
			lines[i] = l[indent:]
		}
	}
	return strings.Join(lines, "\n")
}

var apiGroups = []apiGroup{
	{
		Title: "Auth",
		Endpoints: []apiEndpoint{
			{
				Method: "POST", Path: "/api/waiter/login",
				Summary: "Step 1: verify the restaurant.",
				Params: []apiParam{
					{"restaurant_name", "body", "string", true, "Restaurant name."},
					{"password", "body", "string", true, "The restaurant's password."},
				},
				Request: dedent(`
					{
					  "restaurant_name": "Pasta Place",
					  "password": "restaurant-secret"
					}`),
				Responses: []apiResponse{
					{200, "Restaurant verified.", dedent(`
						{
						  "restaurant_id": 3,
						  "restaurant_name": "Pasta Place",
						  "logo_url": "/uploads/restaurants/3/logo.png"
						}`)},
					errResp(400, "Missing field or broken JSON.", "restaurant_name and password are required"),
					errResp(401, "Wrong name or password.", "incorrect restaurant name or password"),
				},
				Notes: []string{"logo_url is an empty string when the restaurant has no logo."},
			},
			{
				Method: "POST", Path: "/api/waiter/login/employee",
				Summary: "Step 2: verify the employee inside that restaurant and get a bearer token.",
				Params: []apiParam{
					{"restaurant_id", "body", "integer", true, "restaurant_id from step 1."},
					{"employee_name", "body", "string", true, "Employee name."},
					{"password", "body", "string", true, "Employee password."},
				},
				Request: dedent(`
					{
					  "restaurant_id": 3,
					  "employee_name": "Jonas",
					  "password": "employee-secret"
					}`),
				Responses: []apiResponse{
					{200, "Logged in.", dedent(`
						{
						  "token": "eyJhbGciOiJIUzI1NiIs...",
						  "employee_id": 12,
						  "name": "Jonas",
						  "role": "waiter",
						  "expires_at": "2026-11-10T14:30:00+02:00"
						}`)},
					errResp(400, "Missing field or broken JSON.", "restaurant_id, employee_name and password are required"),
					errResp(401, "Wrong name or password.", "incorrect name or password"),
					errResp(403, "The employee is not a waiter or manager.", "the waiter app is for waiters and managers only"),
				},
				Notes: []string{"role is \"waiter\" or \"manager\"; other roles are refused."},
			},
			{
				Method: "GET", Path: "/api/waiter/me", Auth: true,
				Summary: "Who the current token belongs to.",
				Responses: []apiResponse{
					{200, "Identity taken from the token.", dedent(`
						{
						  "employee_id": 12,
						  "name": "Jonas",
						  "role": "waiter",
						  "restaurant_id": 3
						}`)},
					errUnauthorized,
				},
			},
		},
	},
	{
		Title: "Menus & Items",
		Endpoints: []apiEndpoint{
			{
				Method: "GET", Path: "/api/waiter/menus", Auth: true,
				Summary: "All menus of the logged-in restaurant.",
				Responses: []apiResponse{
					{200, "Array of menus.", dedent(`
						[
						  { "id": 5, "name": "Lunch", "description": "Served 11:00-15:00" }
						]`)},
					errUnauthorized,
				},
			},
			{
				Method: "GET", Path: "/api/waiter/menus/{menuID}/items", Auth: true,
				Summary: "Items on one menu.",
				Params: []apiParam{
					{"menuID", "path", "integer", true, "Menu id."},
				},
				Responses: []apiResponse{
					{200, "Array of items.", dedent(`
						[
						  {
						    "id": 41,
						    "name": "Carbonara",
						    "description": "Spaghetti, egg, guanciale",
						    "price": 9.5,
						    "image_url": "/uploads/items/41.jpg",
						    "is_available": true
						  }
						]`)},
					errResp(400, "menuID is not a number.", "invalid menu id"),
					errUnauthorized,
				},
			},
			{
				Method: "GET", Path: "/api/waiter/items/{itemID}", Auth: true,
				Summary: "Full detail of one item.",
				Params: []apiParam{
					{"itemID", "path", "integer", true, "Item id."},
				},
				Responses: []apiResponse{
					{200, "Item detail (English text).", dedent(`
						{
						  "id": 41,
						  "name": "Carbonara",
						  "description": "Spaghetti, egg, guanciale",
						  "instructions": "Boil pasta, ...",
						  "image_url": "/uploads/items/41.jpg",
						  "price": 9.5,
						  "is_available": true,
						  "ingredients": [
						    { "name": "Spaghetti", "measurement": "g", "amount": 120 }
						  ],
						  "allergens": ["Gluten", "Eggs"]
						}`)},
					errResp(400, "itemID is not a number.", "invalid item id"),
					errResp(404, "No such item.", "item not found"),
					errUnauthorized,
				},
			},
			{
				Method: "GET", Path: "/api/waiter/search", Auth: true,
				Summary: "Search the restaurant's items by name.",
				Params: []apiParam{
					{"q", "query", "string", true, "Search text. An empty q returns an empty array."},
				},
				Responses: []apiResponse{
					{200, "Array of items, same shape as the menu item list.", dedent(`
						[
						  {
						    "id": 41,
						    "name": "Carbonara",
						    "description": "Spaghetti, egg, guanciale",
						    "price": 9.5,
						    "image_url": "/uploads/items/41.jpg",
						    "is_available": true
						  }
						]`)},
					errUnauthorized,
				},
			},
		},
	},
	{
		Title: "Quick Items (Pinned)",
		Endpoints: []apiEndpoint{
			{
				Method: "GET", Path: "/api/waiter/quick-items", Auth: true,
				Summary: "The current employee's pinned items.",
				Responses: []apiResponse{
					{200, "Array, ordered by sort_order.", dedent(`
						[
						  {
						    "item_id": 41,
						    "name": "Carbonara",
						    "image_url": "/uploads/items/41.jpg",
						    "price": 9.5,
						    "sort_order": 1
						  }
						]`)},
					errUnauthorized,
				},
			},
			{
				Method: "POST", Path: "/api/waiter/quick-items", Auth: true,
				Summary: "Pin an item.",
				Params: []apiParam{
					{"item_id", "body", "integer", true, "Item to pin."},
					{"sort_order", "body", "integer", false, "Position in the pinned list."},
				},
				Request: `{ "item_id": 41, "sort_order": 1 }`,
				Responses: []apiResponse{
					{201, "Pinned.", okBody},
					errResp(400, "Missing item_id or broken JSON.", "item_id is required"),
					errResp(409, "The employee already has 15 pinned items.", "quick item limit reached (max 15)"),
					errUnauthorized,
				},
			},
			{
				Method: "PATCH", Path: "/api/waiter/quick-items/{itemID}", Auth: true,
				Summary: "Change a pinned item's position.",
				Params: []apiParam{
					{"itemID", "path", "integer", true, "Pinned item id."},
					{"sort_order", "body", "integer", true, "New position."},
				},
				Request: `{ "sort_order": 3 }`,
				Responses: []apiResponse{
					{200, "Updated.", okBody},
					errResp(400, "itemID is not a number or broken JSON.", "invalid item id"),
					errUnauthorized,
				},
			},
			{
				Method: "DELETE", Path: "/api/waiter/quick-items/{itemID}", Auth: true,
				Summary: "Unpin an item.",
				Params: []apiParam{
					{"itemID", "path", "integer", true, "Pinned item id."},
				},
				Responses: []apiResponse{
					{200, "Unpinned.", okBody},
					errResp(400, "itemID is not a number.", "invalid item id"),
					errUnauthorized,
				},
			},
		},
	},
	{
		Title: "Tables",
		Endpoints: []apiEndpoint{
			{
				Method: "GET", Path: "/api/waiter/tables", Auth: true,
				Summary: "All tables of the restaurant.",
				Responses: []apiResponse{
					{200, "Array of tables.", dedent(`
						[
						  { "id": 7, "number": 4, "capacity": 4, "status": "occupied", "has_active_order": true }
						]`)},
					errUnauthorized,
				},
				Notes: []string{
					"status is \"free\", \"occupied\" or \"reserved\".",
					"has_active_order is true while the table has an order that is not delivered or cancelled.",
				},
			},
		},
	},
	{
		Title: "Orders",
		Endpoints: []apiEndpoint{
			{
				Method: "POST", Path: "/api/waiter/orders", Auth: true,
				Summary: "Submit an order. Items are added to the table's open order if there is one.",
				Params: []apiParam{
					{"idempotency_key", "body", "string", true, "Unique per submission, at most 64 characters (the app sends a UUID v4)."},
					{"table_id", "body", "integer | null", false, "null or absent means take-out."},
					{"notes", "body", "string", false, "Order note. Used only when a new order is created."},
					{"items", "body", "array", true, "At least one line."},
					{"items[].item_id", "body", "integer", true, "Item to order."},
					{"items[].quantity", "body", "integer", true, "1 to 1000."},
					{"items[].notes", "body", "string", false, "Note for this line."},
				},
				Request: dedent(`
					{
					  "idempotency_key": "6f1c2c1e-9d0a-4a52-8c1b-0c8d7f3e2b11",
					  "table_id": 7,
					  "notes": "Birthday",
					  "items": [
					    { "item_id": 41, "quantity": 2, "notes": "no cheese" }
					  ]
					}`),
				Responses: []apiResponse{
					{201, "A new order was created.", dedent(`
						{
						  "order_id": 88,
						  "order_number": 15,
						  "status": "pending",
						  "duplicate": false,
						  "appended": false
						}`)},
					{200, "Items were appended to the table's open order (appended: true), or the key was already used and the earlier result is returned (duplicate: true).", dedent(`
						{
						  "order_id": 88,
						  "order_number": 15,
						  "status": "pending",
						  "duplicate": false,
						  "appended": true
						}`)},
					errResp(400, "Invalid body (also: key longer than 64, empty items, quantity outside 1-1000).", "quantity must be between 1 and 1000"),
					errUnauthorized,
				},
				Notes: []string{
					"Take-out orders (no table_id) are never merged: each one is a new order.",
					"Safe to retry with the same idempotency_key: it will not create a second order or add the items twice.",
					"order_number restarts every day at 00:00.",
				},
			},
			{
				Method: "GET", Path: "/api/waiter/orders", Auth: true,
				Summary: "The current employee's own active orders.",
				Responses: []apiResponse{
					{200, "Array of orders that are not delivered or cancelled.", dedent(`
						[
						  {
						    "order_id": 88,
						    "order_number": 15,
						    "table_number": 4,
						    "status": "pending",
						    "created_at": "2026-11-10 13:05",
						    "notes": "Birthday",
						    "items": [
						      { "item_id": 41, "name": "Carbonara", "quantity": 2, "price": 9.5 }
						    ]
						  }
						]`)},
					errUnauthorized,
				},
				Notes: []string{
					"table_number is null for take-out.",
					"status is \"pending\", \"preparing\" or \"ready\" here.",
					"name and price are the values at the time of ordering.",
				},
			},
			{
				Method: "PATCH", Path: "/api/waiter/orders/{orderID}/items/{itemOrderID}", Auth: true,
				Summary: "Change the quantity of one order line.",
				Params: []apiParam{
					{"orderID", "path", "integer", true, "Order id."},
					{"itemOrderID", "path", "integer", true, "Order line id."},
					{"quantity", "body", "integer", true, "1 to 1000."},
				},
				Request: `{ "quantity": 3 }`,
				Responses: []apiResponse{
					{200, "Updated.", okBody},
					errResp(400, "Bad ids or quantity.", "quantity must be between 1 and 1000"),
					errResp(404, "No such order.", "order not found"),
					errResp(409, "The order is no longer pending.", "order is already being prepared and can no longer be edited"),
					errUnauthorized,
				},
				Notes: []string{"Only allowed while the order is pending."},
			},
			{
				Method: "DELETE", Path: "/api/waiter/orders/{orderID}/items/{itemOrderID}", Auth: true,
				Summary: "Remove one order line.",
				Params: []apiParam{
					{"orderID", "path", "integer", true, "Order id."},
					{"itemOrderID", "path", "integer", true, "Order line id."},
				},
				Responses: []apiResponse{
					{200, "Removed.", okBody},
					errResp(400, "Bad ids.", "invalid order or item id"),
					errResp(404, "No such order.", "order not found"),
					errResp(409, "The order is no longer pending.", "order is already being prepared and can no longer be edited"),
					errUnauthorized,
				},
				Notes: []string{"Only allowed while the order is pending."},
			},
			{
				Method: "POST", Path: "/api/waiter/orders/{orderID}/cancel", Auth: true,
				Summary: "Cancel an order.",
				Params: []apiParam{
					{"orderID", "path", "integer", true, "Order id."},
				},
				Responses: []apiResponse{
					{200, "Cancelled.", okBody},
					errResp(400, "orderID is not a number.", "invalid order id"),
					errResp(404, "No such order.", "order not found"),
					errResp(409, "The order is no longer pending.", "order is already being prepared and can no longer be cancelled"),
					errUnauthorized,
				},
				Notes: []string{"Only allowed while the order is pending."},
			},
		},
	},
	{
		Title: "Chat",
		Endpoints: []apiEndpoint{
			{
				Method: "GET", Path: "/api/waiter/chat/messages", Auth: true,
				Summary: "All chat messages of the restaurant.",
				Responses: []apiResponse{
					{200, "Array of messages.", dedent(`
						[
						  { "id": 120, "sender_name": "Jonas", "message": "Table 4 needs water", "sent_at": "2026-11-10 13:10" }
						]`)},
					errUnauthorized,
				},
			},
			{
				Method: "GET", Path: "/api/waiter/chat/messages/after/{id}", Auth: true,
				Summary: "Messages newer than the given message id. Use it for polling.",
				Params: []apiParam{
					{"id", "path", "integer", true, "Id of the newest message the client already has."},
				},
				Responses: []apiResponse{
					{200, "Array of newer messages, same shape as above (empty if none).", dedent(`
						[
						  { "id": 121, "sender_name": "Ona", "message": "On it", "sent_at": "2026-11-10 13:11" }
						]`)},
					errResp(400, "id is not a number.", "invalid id"),
					errUnauthorized,
				},
			},
			{
				Method: "POST", Path: "/api/waiter/chat/send", Auth: true,
				Summary: "Send a chat message.",
				Params: []apiParam{
					{"message", "body", "string", true, "Message text."},
				},
				Request: `{ "message": "Table 4 needs water" }`,
				Responses: []apiResponse{
					{201, "Sent.", okBody},
					errResp(400, "Empty message or broken JSON.", "message is required"),
					errUnauthorized,
				},
			},
		},
	},
}
