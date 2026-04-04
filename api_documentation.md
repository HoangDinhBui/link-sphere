# LinkSphere API Documentation

This document covers all public API endpoints exposed by the LinkSphere OpenShift API Gateway. 

**Base URLs**:
*   Local cluster (Docker-Compose): `http://localhost:8888`
*   OpenShift Route: `http://api-gateway-buidinhhoang1910-dev.apps.rm3.7wse.p1.openshiftapps.com`

**Global Headers**:
For protected endpoints (marked with 🔒), you must include the following header:
`Authorization: Bearer <your_access_token>`

All JSON Payload requests should include:
`Content-Type: application/json`

---

## 1. Authentication Service (`/api/v1/auth`)

### 1.1 Login
*   **Method**: `POST`
*   **Path**: `/api/v1/auth/login`
*   **Description**: Authenticates user and returns JWT token.
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "secretpassword"
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "success": true,
      "message": "login successful",
      "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR..."
      }
    }
    ```

---

## 2. User Service (`/api/v1/users`)

### 2.1 Register User
*   **Method**: `POST`
*   **Path**: `/api/v1/users/register`
*   **Description**: Creates a new user account.
*   **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "username": "user123",
      "password": "secretpassword",
      "phone": "0123456789"  // Optional
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "success": true,
      "message": "user registered successfully",
      "data": {
        "id": "uuid-here",
        "email": "user@example.com",
        "username": "user123",
        "created_at": "2024-04-10T12:00:00Z"
      }
    }
    ```

### 2.2 Get Current Profile 🔒
*   **Method**: `GET`
*   **Path**: `/api/v1/users/profile`
*   **Response (200 OK)**:
    ```json
    {
      "id": "uuid-here",
      "email": "...",
      "username": "..."
    }
    ```

### 2.3 Follow a User 🔒
*   **Method**: `POST`
*   **Path**: `/api/v1/users/follow`
*   **Request Body**:
    ```json
    {
      "targetUserId": "uuid-of-user-to-follow"
    }
    ```

### 2.4 Unfollow a User 🔒
*   **Method**: `POST`
*   **Path**: `/api/v1/users/unfollow`
*   **Request Body**:
    ```json
    {
      "targetUserId": "uuid-to-unfollow"
    }
    ```

### 2.5 Inspect Followings 🔒
*   **Method**: `GET`
*   **Path**: `/api/v1/users/{id}/following`
*   **Response (200 OK)**: array of UUIDs.

---

## 3. Post Service (`/api/v1/posts`)

### 3.1 Create Post 🔒
*   **Method**: `POST`
*   **Path**: `/api/v1/posts`
*   **Request Body**:
    ```json
    {
      "content": "Hello LinkSphere!",
      "images": ["https://s3/path/img1.png"], // Optional
      "hashtags": ["hello", "world"]          // Optional
    }
    ```
*   **Response (200 OK)**:
    ```json
    {
      "success": true,
      "message": "post created successfully",
      "data": {
        "id": "uuid",
        "user_id": "uuid",
        "content": "Hello LinkSphere!",
        "like_count": 0,
        "created_at": "2024-04-10T..."
      }
    }
    ```

### 3.2 List Global Posts 
*   **Method**: `GET`
*   **Path**: `/api/v1/posts`
*   **Body** (sent via GET Body/Params depending on standard client behavior):
    ```json
    {
       "page": 1,
       "limit": 10
    }
    ```

### 3.3 Get Posts by Specific Users (Bulk Fetch)
*   **Method**: `POST`
*   **Path**: `/api/v1/posts/by-users`
*   **Request Body**:
    ```json
    {
      "userIds": ["uuid1", "uuid2"],
      "page": 1,
      "limit": 20
    }
    ```

### 3.4 Like a Post 🔒
*   **Method**: `POST`
*   **Path**: `/api/v1/posts/{post_id}/like` *(Assumed URL format, depends on chi params)*
*   **Response (200 OK)**:
    ```json
    { "success": true, "message": "post liked successfully" }
    ```

---

## 4. Comment Service (`/api/v1/posts/comment`)

### 4.1 Create Comment 🔒
*   **Method**: `POST`
*   **Path**: `/api/v1/posts/comment`
*   **Request Body**:
    ```json
    {
      "postId": "uuid-of-the-post",
      "content": "This is a great comment!"
    }
    ```

### 4.2 List Comments
*   **Method**: `GET`
*   **Path**: `/api/v1/posts/comments`
*   **Request Body**:
    ```json
    {
      "postId": "uuid-of-the-post",
      "page": 1,
      "limit": 20
    }
    ```

---

## 5. Feed Service (`/api/v1/feed`)

### 5.1 Get Personalized News Feed 🔒
*   **Method**: `GET`
*   **Path**: `/api/v1/feed`
*   **Request Body**:
    ```json
    {
      "page": 1,
      "limit": 10
    }
    ```
*   **Description**: Aggregates posts from the people you follow. Relies on internal cache (Redis) and sub-queries to Post & User services.

---

## 6. Search Service (`/api/v1/search`)

### 6.1 Search Resources
*   **Method**: `GET`
*   **Path**: 
    - `/api/v1/search/posts?q={keyword}`
    - `/api/v1/search/users?q={keyword}`

---

## 7. Notification Service (`/api/v1/notifications`)

### 7.1 WebSocket Connection 🔒
*   **Endpoint**: `ws://api-gateway-...openshiftapps.com/ws/`
*   **Description**: Maintains a persistent connection to receive realtime Socket events (Likes, Comments).

### 7.2 Get Past Notifications 🔒
*   **Method**: `GET`
*   **Path**: `/api/v1/notifications`
