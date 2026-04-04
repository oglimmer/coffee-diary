import AuthenticationServices
import SwiftUI

@Observable
final class AuthViewModel {
    var user: User?
    var isLoading = true
    var error: String?

    var isAuthenticated: Bool { user != nil }

    private let authService = AuthService()

    func checkSession() async {
        isLoading = true
        user = await authService.checkSession()
        isLoading = false
    }

    func login(anchor: ASPresentationAnchor) async {
        error = nil
        do {
            try await authService.login(anchor: anchor)
            user = await authService.checkSession()
        } catch {
            self.error = error.localizedDescription
        }
    }

    func loginWithApple(result: Result<ASAuthorization, any Error>) async {
        error = nil
        do {
            guard case .success(let authorization) = result,
                  let credential = authorization.credential as? ASAuthorizationAppleIDCredential,
                  let identityToken = credential.identityToken else {
                self.error = "Apple Sign-In failed"
                return
            }

            var fullName: String?
            if let nameComponents = credential.fullName {
                let parts = [nameComponents.givenName, nameComponents.familyName].compactMap { $0 }
                if !parts.isEmpty {
                    fullName = parts.joined(separator: " ")
                }
            }

            try await authService.loginWithApple(identityToken: identityToken, fullName: fullName)
            user = await authService.checkSession()
        } catch {
            self.error = error.localizedDescription
        }
    }

    func logout() async {
        await authService.logout()
        user = nil
    }
}
