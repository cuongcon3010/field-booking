# Field Booking - Google OAuth 2.0 Authentication

Hệ thống xác thực Google OAuth 2.0 an toàn cho ứng dụng Field Booking, phát triển bằng **Go (Gin, GORM, PostgreSQL)** và **Next.js (React, Tailwind CSS)**.

---

## 1. Tính năng chính

- **Đăng nhập & Đăng ký qua Google**: Tự động tạo user mới nếu chưa tồn tại.
- **Tự động liên kết tài khoản (Account Linking)**:
  - Nếu Google ID đã liên kết -> Đăng nhập vào User tương ứng.
  - Nếu Email Google đã tồn tại trong hệ thống và đã được xác thực (`email_verified = true`) -> Tự động tạo bản ghi liên kết `oauth_accounts` vào user hiện tại.
  - Từ chối đăng nhập nếu email chưa được xác thực từ Google (`email_verified = false`).
- **CSRF Protection với OAuth State**:
  - Sinh chuỗi ngẫu nhiên bằng `crypto/rand` (cryptographically secure).
  - Lưu vào HttpOnly Cookie với thời hạn hết hạn ngắn (10 phút).
  - Xác thực so khớp giữa cookie và query param khi Google callback; xóa cookie ngay sau khi xác thực để chống tấn công replay.
- **Bảo mật Cookie & Phiên đăng nhập**:
  - Ứng dụng phát hành JWT token và lưu trữ trong **HttpOnly Cookie** (`SameSite=Lax`, `Path=/`).
  - **Tuyệt đối không lưu token vào localStorage**, ngăn ngừa rủi ro bị đánh cắp qua tấn công XSS.
  - Không bao giờ để lộ Google access token, refresh token hay client secret ra frontend.
- **CORS chuẩn mực với Credentials**:
  - Hỗ trợ `Access-Control-Allow-Credentials: true`.
  - Giới hạn chính xác domain `FRONTEND_URL` (không dùng wildcard `*`).

---

## 2. Thiết kế Database

Sử dụng quan hệ 1 - N giữa `users` và `oauth_accounts` để dễ dàng mở rộng thêm các provider khác sau này (Facebook, Github, Discord...):

```sql
-- Bảng users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- Bảng oauth_accounts
CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT idx_provider_user_id UNIQUE (provider, provider_user_id)
);
```

---

## 3. Hướng dẫn cấu hình Google Cloud Console

1. Truy cập [Google Cloud Console](https://console.cloud.google.com/).
2. Tạo một project mới (ví dụ: `field-booking`).
3. Vào menu **APIs & Services** > **OAuth consent screen**:
   - Chọn **External** (hoặc Internal nếu dùng workspace tổ chức).
   - Điền tên App, Email hỗ trợ.
   - Tại mục **Scopes**, thêm 3 scopes cơ bản:
     - `.../auth/userinfo.email`
     - `.../auth/userinfo.profile`
     - `openid`
4. Vào menu **Credentials** > **Create Credentials** > **OAuth client ID**:
   - Application type: **Web application**.
   - Name: `Field Booking Web Client`.
   - **Authorized JavaScript origins**:
     - `http://localhost:3000`
   - **Authorized redirect URIs**:
     - `http://localhost:8080/api/auth/google/callback`
5. Nhấn **Create** và sao chép **Client ID** cùng **Client Secret**.

---

## 4. Thiết lập biến môi trường

Tạo file `.env` tại thư mục gốc hoặc trong `BackEnd/.env` dựa trên `.env.example`:

```env
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000

# PostgreSQL Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=field_booking
DB_SSLMODE=disable

# JWT Secret
JWT_SECRET=thay-doi-chuoi-nay-trong-moi-truong-that-rat-dai-va-ngau-nhien
JWT_EXPIRATION_HOURS=72

# Google OAuth 2.0
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

---

## 5. Hướng dẫn khởi chạy

### Bước 1: Khởi động Database (PostgreSQL)

* **Sử dụng Podman:**
  ```bash
### Cách 1: Chạy trọn bộ hệ thống bằng Compose (Khuyên dùng)

Khởi động đồng thời cả 3 dịch vụ (**PostgreSQL**, **Go Backend**, **Next.js Frontend**):

* **Sử dụng Podman:**
  ```bash
  podman compose up -d
  # hoặc: podman-compose up -d
  ```

* **Sử dụng Docker:**
  ```bash
  docker compose up -d
  ```

Truy cập giao diện tại: `http://localhost:3000` (Backend API tại `http://localhost:8080`).

---

### Cách 2: Chạy thủ công từng thành phần cho môi trường phát triển

1. **Chạy riêng Database:**
   ```bash
   podman compose up -d postgres
   ```

2. **Chạy Backend (Go):**
   ```bash
   cd BackEnd
   go run ./cmd/server
   ```

3. **Chạy Frontend (Next.js):**
   ```bash
   cd FrontEnd
   npm run dev
   ```

---

## 6. Danh sách API Endpoints

| Method | Endpoint | Mô tả | Yêu cầu Auth |
|---|---|---|---|
| `GET` | `/api/auth/google` | Khởi tạo OAuth state, set cookie CSRF và redirect tới Google | Không |
| `GET` | `/api/auth/google/callback` | Nhận auth code, xác minh CSRF, đăng nhập/tạo user, set HttpOnly cookie | Không |
| `GET` | `/api/auth/me` | Lấy thông tin user hiện tại qua HttpOnly cookie | Có (`access_token` cookie) |
| `POST` | `/api/auth/logout` | Đăng xuất, xóa HttpOnly cookie | Không |
| `GET` | `/health` | Healthcheck server backend | Không |

---

## 7. Chạy Tests & Kiểm tra mã nguồn

```bash
cd BackEnd
# Chạy toàn bộ Unit Tests
go test -v ./...

# Kiểm tra cú pháp chuẩn Go
go vet ./...
```

