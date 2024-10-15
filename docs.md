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

