package hub

var externalApiGroups = []apiGroup{
	{
		Title: "Items",
		Endpoints: []apiEndpoint{
			{
				Method: "GET", Path: "/api/external/items", Auth: true,
				Summary: "Items the restaurant currently offers. Use the ids when submitting an order.",
				Responses: []apiResponse{
					{200, "Array of available items.", dedent(`
						[
						  {
						    "id": 41,
						    "name": "Carbonara",
						    "description": "Spaghetti, egg, guanciale",
						    "price": 9.5,
						    "image_url": "/uploads/items/41.jpg"
						  }
						]`)},
					errResp(401, "Missing, invalid or revoked api key.", "invalid or revoked api key"),
				},
				Notes: []string{
					"image_url is a relative path on this CommonTable server.",
					"Only items that are currently available are listed.",
				},
			},
		},
	},
	{
		Title: "Orders",
		Endpoints: []apiEndpoint{
			{
				Method: "POST", Path: "/api/external/orders", Auth: true,
				Summary: "Submit a take-out order. It waits for a manager to approve it.",
				Params: []apiParam{
					{"idempotency_key", "body", "string", true, "Unique per submission, at most 64 characters (a UUID v4 works)."},
					{"customer_name", "body", "string", true, "At most 255 characters."},
					{"customer_phone", "body", "string", true, "At most 30 characters: digits, +, -, ( ) and spaces."},
					{"notes", "body", "string", false, "At most 500 characters."},
					{"items", "body", "array", true, "1 to 50 lines."},
					{"items[].item_id", "body", "integer", true, "An id from the items endpoint."},
					{"items[].quantity", "body", "integer", true, "1 to 1000."},
					{"items[].notes", "body", "string", false, "At most 500 characters."},
				},
				Request: dedent(`
					{
					  "idempotency_key": "6f1c2c1e-9d0a-4a52-8c1b-0c8d7f3e2b11",
					  "customer_name": "Jonas Jonaitis",
					  "customer_phone": "+370 600 00000",
					  "notes": "Pick-up at 18:00",
					  "items": [
					    { "item_id": 41, "quantity": 2, "notes": "no cheese" }
					  ]
					}`),
				Responses: []apiResponse{
					{201, "Order received and waiting for approval.", dedent(`
						{
						  "order_id": 17,
						  "status": "pending",
						  "duplicate": false
						}`)},
					{200, "The idempotency_key was already used. The earlier order is returned and nothing new is created.", dedent(`
						{
						  "order_id": 17,
						  "status": "pending",
						  "duplicate": true
						}`)},
					errResp(400, "Missing or invalid field, or an item that does not exist.", "customer_name and customer_phone are required"),
					errResp(401, "Missing, invalid or revoked api key.", "invalid or revoked api key"),
					errResp(409, "An item is not available right now.", "item 41 is not available"),
					errResp(429, "A limit was reached. A Retry-After header is sent when the cause is the per-minute limit.", "too many orders in a short time, try again later"),
				},
				Notes: []string{
					"Prices and item names are set by the server. Anything the client sends about them is ignored.",
					"The order_id is the id of the waiting order. It is not the number shown in the restaurant.",
					"Safe to retry with the same idempotency_key: it will not create a second order.",
					"Limits per key: 10 orders per minute and 50 waiting orders. Limit per phone number: 3 waiting orders.",
				},
			},
		},
	},
}
