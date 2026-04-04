import Foundation

enum APIError: Error, LocalizedError {
    case unauthorized
    case badRequest(String)
    case serverError(Int)
    case decodingError(Error)
    case networkError(Error)

    var errorDescription: String? {
        switch self {
        case .unauthorized: "Authentication required"
        case .badRequest(let msg): msg
        case .serverError(let code): "Server error (\(code))"
        case .decodingError(let err): "Decoding error: \(err.localizedDescription)"
        case .networkError(let err): err.localizedDescription
        }
    }
}

@Observable
final class APIClient {
    static let shared = APIClient()

    let baseURL = URL(string: "https://coffee.oglimmer.com")!

    private let session: URLSession
    private let decoder: JSONDecoder = {
        let d = JSONDecoder()
        return d
    }()

    private init() {
        let config = URLSessionConfiguration.default
        config.httpCookieAcceptPolicy = .always
        config.httpCookieStorage = .shared
        session = URLSession(configuration: config)
    }

    /// Injects the session cookie received from the OIDC redirect into the shared cookie storage.
    func setSessionCookie(_ value: String) {
        guard let cookie = HTTPCookie(properties: [
            .name: "session",
            .value: value,
            .domain: baseURL.host()!,
            .path: "/",
            .secure: "TRUE",
            .expires: Date().addingTimeInterval(7 * 24 * 3600),
        ]) else { return }
        HTTPCookieStorage.shared.setCookie(cookie)
    }

    func clearSessionCookie() {
        if let cookies = HTTPCookieStorage.shared.cookies(for: baseURL) {
            for cookie in cookies where cookie.name == "session" {
                HTTPCookieStorage.shared.deleteCookie(cookie)
            }
        }
    }

    var hasSessionCookie: Bool {
        HTTPCookieStorage.shared.cookies(for: baseURL)?.contains { $0.name == "session" } ?? false
    }

    // MARK: - Generic request

    func request<T: Decodable>(_ method: String, path: String, queryItems: [URLQueryItem] = [], body: (any Encodable)? = nil) async throws -> T {
        var url = baseURL.appending(path: path)
        if !queryItems.isEmpty {
            url.append(queryItems: queryItems)
        }

        var req = URLRequest(url: url)
        req.httpMethod = method

        if let body {
            req.httpBody = try JSONEncoder().encode(body)
            req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        }

        let (data, response): (Data, URLResponse)
        do {
            (data, response) = try await session.data(for: req)
        } catch {
            throw APIError.networkError(error)
        }

        guard let http = response as? HTTPURLResponse else {
            throw APIError.serverError(0)
        }

        switch http.statusCode {
        case 200...201:
            do {
                return try decoder.decode(T.self, from: data)
            } catch {
                throw APIError.decodingError(error)
            }
        case 401:
            throw APIError.unauthorized
        case 400...499:
            let msg = (try? decoder.decode(ErrorResponse.self, from: data))?.message ?? "Bad request"
            throw APIError.badRequest(msg)
        default:
            throw APIError.serverError(http.statusCode)
        }
    }

    func requestNoContent(_ method: String, path: String) async throws {
        let url = baseURL.appending(path: path)
        var req = URLRequest(url: url)
        req.httpMethod = method

        let (data, response): (Data, URLResponse)
        do {
            (data, response) = try await session.data(for: req)
        } catch {
            throw APIError.networkError(error)
        }

        guard let http = response as? HTTPURLResponse else {
            throw APIError.serverError(0)
        }

        if http.statusCode == 401 { throw APIError.unauthorized }
        if http.statusCode >= 400 {
            let msg = (try? decoder.decode(ErrorResponse.self, from: data))?.message ?? "Request failed"
            throw APIError.badRequest(msg)
        }
    }
}

private struct ErrorResponse: Decodable {
    let message: String
}
