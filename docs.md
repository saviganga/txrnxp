# API Documentation


## APIs Overview

- **Base Url**: `127.0.0.1:8000`

### Authentication
The application uses JWT Authentication for secure access to specific endpoints

#### 1. Signup
This endpoint allows users to create a new account.

- **Endpoint**: `{{BASE_URL}}/api/v1/users/`
- **Method**: POST

**Request Body:**
```json
{
   "email": "ganga2@ganga.com",
    "password": "ganga1234",
    "username": "ganga2",
    "first_name": "ganga2",
    "last_name": "ganga2"
}
```

#### 2. Login
This endpoint allows users to login to their accounts.

- **Endpoint**: `{{BASE_URL}}/api/v1/auth/login`
- **Method**: POST

**Request Body:**
```json
{
    "username": "jondoe",
    "password": "xxxxxx"
}
```

**Request Headers:**
```json
{
    "platform": "USER"
}
```


#### 3. Get user data
This endpoint allows users to view their user account information.

- **Endpoint**: `{{BASE_URL}}/api/v1/users/`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 4. Retrieve user data
This endpoint allows users to retrieve their user account information.

- **Endpoint**: `{{BASE_URL}}/api/v1/users/<user_id>`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 5. Update user data
This endpoint allows users to update their user account information.

- **Endpoint**: `{{BASE_URL}}/api/v1/users/<user_id>`
- **Method**: PATCH
- **Authorization**: true

**Request Body:**
```json
{
    "username": "jondoe",
    "password": "xxxxxx"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 6. Upload user image
This endpoint allows users to upload a profile picture.

- **Endpoint**: `{{BASE_URL}}/api/v1/users/<user_id>/upload-image`
- **Method**: POST
- **Authorization**: true

**Request Body:**
```json
{
    "image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


### Wallets

#### 1. Get user wallet
This endpoint allows users to view their wallet information.

- **Endpoint**: `{{BASE_URL}}/api/v1/wallets`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 2. Get user wallet entries
This endpoint allows users to view their wallet transactions.

- **Endpoint**: `{{BASE_URL}}/api/v1/wallets/entries`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 3. Transfer
This endpoint allows users to transfer funds between wallets.

- **Endpoint**: `{{BASE_URL}}/api/v1/wallets/transfer`
- **Method**: POST
- **Authorization**: true


**Request body:**
```json
{
    "receiver_email": "ganga2@ganga.com",
    "amount": "500"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

### Business


#### 1. Create Business
This endpoint allows users to create businesses.

- **Endpoint**: `{{BASE_URL}}/api/v1/business/`
- **Method**: POST

**Request Body:**
```json
{
    "name": "ganga2-1"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

#### 2. Get Businesses
This endpoint allows users to view their businesses.

- **Endpoint**: `{{BASE_URL}}/api/v1/business/`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 3. Get Businesses
This endpoint allows users to view their business.

- **Endpoint**: `{{BASE_URL}}/api/v1/business/<business_id>`
- **Method**: GET
- **Authorization**: true

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 4. Update Businesses
This endpoint allows users to update their business.

- **Endpoint**: `{{BASE_URL}}/api/v1/business/<business_id>`
- **Method**: PATCH
- **Authorization**: true


**Request Body:**
```json
{
    "name": "ganga2-1"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


#### 5. Upload business image
This endpoint allows businesses to upload a profile picture.

- **Endpoint**: `{{BASE_URL}}/api/v1/businesses/<business_id>/upload-image`
- **Method**: POST
- **Authorization**: true

**Request Body:**
```json
{
    "image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/"
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```




### Events


#### 1. Create Event
This endpoint allows users to create events. Events can be created by individuals or businesses

- **Endpoint**: `{{BASE_URL}}/api/v1/events/`
- **Method**: POST

**Request Body:**
```json
{
   "name": "lagos nights",
    "category": "PARTY",
    "duration": "recurring",
    "event_type": "LIVE",
    "address": "victoria island lagos",
    "description": "beach party",
    "start_time": "2024-11-06T13:38:43.827061Z",
    "end_time": "2024-12-06T13:38:43.827061Z"
}
```

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```


#### 2. Upload Event Image
This endpoint allows users to upload images for their created events.

- **Endpoint**: `{{BASE_URL}}/api/v1/events/<event_id>/upload-image`
- **Method**: POST
- **Authorization**: true


**Request Body:**
```json
{
    "image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/"
}
```

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```


**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```



#### 3. Get My Events
This endpoint allows users to view their created events.

- **Endpoint**: `{{BASE_URL}}/api/v1/events/my-events`
- **Method**: GET
- **Authorization**: true

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```


#### 4. Update My Event
This endpoint allows users to update their created events.

- **Endpoint**: `{{BASE_URL}}/api/v1/events/<event_id>`
- **Method**: PATCH
- **Authorization**: true


**Request Body:**
```json
{
    "name": "lagos vibez"
}
```


**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```


#### 3. Get Upcoming Events
This endpoint allows users to view upcoming events.

- **Endpoint**: `{{BASE_URL}}/api/v1/events`
- **Method**: GET
- **Authorization**: false


#### 4. Get Event Detail
This endpoint allows users to view event details.

- **Endpoint**: `{{BASE_URL}}/api/v1/events/<event_reference>`
- **Method**: GET
- **Authorization**: false


#### 4. Get Event History
This endpoint allows users to view events they attended/purchased tickets for. This feature is not available for businesses

- **Endpoint**: `{{BASE_URL}}/api/v1/events/history`
- **Method**: GET
- **Authorization**: true

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "USER"
}
```


#### 5. Create Event Ticket
This endpoint allows users to create event tickets.

- **Endpoint**: `{{BASE_URL}}/api/v1/tickets/events/`
- **Method**: POST

**Request Body:**
```json
{
   "event_id": "d83c6d7d-039e-499a-b9be-75d243955bbb",
    "is_paid": true,
    "is_limited_stock": true,
    "price": "750",
    "perks": {
        "food": false,
        "drinks": true,
        "vip": false
    },
    "stock_number": 15,
    "is_invite_only": false,
    "ticket_type": "SINGLE",
    "description": "event ticket",
    "purchase_limit": 3,
    "is_limited_stock": true,
    "stock_number": 100
}
```

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```



#### 6. Get Event Tickets
This endpoint allows users to view event tickets. Only event organisers have access to the feature

- **Endpoint**: `{{BASE_URL}}/api/v1/events/<event_id>/tickets`
- **Method**: GET
- **Authorization**: true

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}",
}
```

**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```


#### 7. Get Event Ticket by Id
This endpoint allows users to view event tickets. Only event organisers have access to the feature

- **Endpoint**: `{{BASE_URL}}/api/v1/tickets/events/<event_ticket_id>`
- **Method**: GET
- **Authorization**: true

**Request Headers: Individual**
```json
{
    "Authorization": "JWT {{jwt_token}}",
}
```

**Request Headers: Business**
```json
{
    "Authorization": "JWT {{jwt_token}}",
    "Entity": "BUSINESS",
    "Business": "<business_reference>"
}
```


### User Tickets


#### 1. Buy Ticket
This endpoint allows users to buy event tickets

- **Endpoint**: `{{BASE_URL}}/api/v1/tickets/buy/wallet`
- **Method**: POST

**Request Body:**
```json
{
   "event_ticket_id": "317ffc72-68ba-4a88-88e2-c78413246de0",
    "count": 2
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

#### 2. Transfer Ticket
This endpoint allows users to transfer event tickets

- **Endpoint**: `{{BASE_URL}}/api/v1/tickets/transfer`
- **Method**: POST

**Request Body:**
```json
{
   "ticket_reference": "<ticket_reference>",
   "receiver_email": "nigga@niggas.com",
    "count": 2
}
```

**Request Headers:**
```json
{
    "Authorization": "JWT {{jwt_token}}"
}
```

