import AuthenticationServices
import SwiftUI

struct LoginError {
    let message: String
    let detail: String?
}

@Observable
final class AuthViewModel {
    var user: User?
    var isLoading = true
    var error: LoginError?

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
        } catch let error as ASWebAuthenticationSessionError where error.code == .canceledLogin {
            // User cancelled — not an error, just do nothing
        } catch {
            self.error = LoginError(
                message: "Sign-in failed. Please check your connection and try again.",
                detail: error.localizedDescription
            )
        }
    }

    func loginWithApple(result: Result<ASAuthorization, any Error>) async {
        error = nil
        do {
            guard case .success(let authorization) = result,
                  let credential = authorization.credential as? ASAuthorizationAppleIDCredential,
                  let identityToken = credential.identityToken else {
                if case .failure(let resultError) = result {
                    let nsError = resultError as NSError
                    if nsError.domain == ASAuthorizationError.errorDomain,
                       nsError.code == ASAuthorizationError.canceled.rawValue {
                        return // User cancelled — not an error
                    }
                }
                self.error = LoginError(
                    message: "Apple Sign-In failed. Please try again.",
                    detail: nil
                )
                return
            }

            var fullName: String?
            if let nameComponents = credential.fullName {
                let parts = [nameComponents.givenName, nameComponents.familyName].compactMap { $0 }
                if !parts.isEmpty {
                    fullName = parts.joined(separator: " ")
                }
            }

            try await authService.loginWithApple(
                identityToken: identityToken,
                authorizationCode: credential.authorizationCode,
                fullName: fullName
            )
            user = await authService.checkSession()
        } catch {
            self.error = LoginError(
                message: "Apple Sign-In failed. Please check your connection and try again.",
                detail: error.localizedDescription
            )
        }
    }

    func logout() async {
        await authService.logout()
        user = nil
    }

    func deleteAccount() async {
        error = nil
        do {
            try await authService.deleteAccount()
            user = nil
        } catch {
            self.error = LoginError(
                message: "Account deletion failed. Please try again.",
                detail: error.localizedDescription
            )
        }
    }
}
