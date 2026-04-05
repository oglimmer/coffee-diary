import AuthenticationServices
import Foundation

final class AuthService {
    private let api = APIClient.shared

    /// Retained so ASWebAuthenticationSession isn't deallocated mid-flow.
    private var webAuthSession: ASWebAuthenticationSession?
    private var contextProvider: PresentationContextProvider?

    /// Opens the OIDC login flow in an in-app browser and returns the session cookie.
    func login(anchor: ASPresentationAnchor) async throws {
        let redirectScheme = "coffeeDiary"
        let loginURL = api.baseURL
            .appending(path: "/api/auth/login")
            .appending(queryItems: [
                .init(name: "redirect_after", value: "coffeeDiary://auth/callback"),
            ])

        let callbackURL: URL = try await withCheckedThrowingContinuation { continuation in
            let session = ASWebAuthenticationSession(
                url: loginURL,
                callbackURLScheme: redirectScheme
            ) { url, error in
                if let error {
                    continuation.resume(throwing: error)
                } else if let url {
                    continuation.resume(returning: url)
                } else {
                    continuation.resume(throwing: APIError.unauthorized)
                }
            }
            let provider = PresentationContextProvider(anchor: anchor)
            session.presentationContextProvider = provider
            session.prefersEphemeralWebBrowserSession = true

            // Retain both until the completion handler fires
            self.webAuthSession = session
            self.contextProvider = provider

            session.start()
        }

        // Clean up
        webAuthSession = nil
        contextProvider = nil

        // Extract session cookie from callback URL
        guard let components = URLComponents(url: callbackURL, resolvingAgainstBaseURL: false),
              let cookieValue = components.queryItems?.first(where: { $0.name == "session_cookie" })?.value,
              !cookieValue.isEmpty
        else {
            throw APIError.unauthorized
        }

        api.setSessionCookie(cookieValue)
    }

    /// Sends the Apple identity token and authorization code to the backend for verification
    /// and session creation. The authorization code is required so the backend can obtain a
    /// refresh token and revoke it at account deletion (App Store Guideline 5.1.1(v)).
    func loginWithApple(identityToken: Data, authorizationCode: Data?, fullName: String?) async throws {
        guard let tokenString = String(data: identityToken, encoding: .utf8) else {
            throw APIError.badRequest("Invalid identity token")
        }
        let codeString = authorizationCode.flatMap { String(data: $0, encoding: .utf8) }

        struct AppleLoginRequest: Encodable {
            let identityToken: String
            let authorizationCode: String?
            let fullName: String?
        }

        let _: User = try await api.request(
            "POST",
            path: "/api/auth/apple-callback",
            body: AppleLoginRequest(
                identityToken: tokenString,
                authorizationCode: codeString,
                fullName: fullName
            )
        )
    }

    /// Permanently deletes the authenticated user's account and all their data.
    func deleteAccount() async throws {
        try await api.requestNoContent("DELETE", path: "/api/auth/me")
        api.clearSessionCookie()
    }

    func checkSession() async -> User? {
        try? await api.request("GET", path: "/api/auth/me")
    }

    func logout() async {
        try? await api.requestNoContent("GET", path: "/api/auth/logout")
        api.clearSessionCookie()
    }
}

// MARK: - ASWebAuthenticationSession presentation

private final class PresentationContextProvider: NSObject, ASWebAuthenticationPresentationContextProviding, @unchecked Sendable {
    let anchor: ASPresentationAnchor
    init(anchor: ASPresentationAnchor) { self.anchor = anchor }
    func presentationAnchor(for _: ASWebAuthenticationSession) -> ASPresentationAnchor { anchor }
}
