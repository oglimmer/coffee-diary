import AuthenticationServices
import SwiftUI

struct LoginView: View {
    var authViewModel: AuthViewModel
    @Environment(\.colorScheme) private var colorScheme
    @State private var window: UIWindow?

    var body: some View {
        VStack(spacing: 32) {
            Spacer()

            Image(systemName: "cup.and.saucer.fill")
                .font(.system(size: 80))
                .foregroundStyle(.brown)

            Text("Coffee Diary")
                .font(.largeTitle.bold())

            Text("Track your espresso brewing sessions")
                .font(.subheadline)
                .foregroundStyle(.secondary)

            Spacer()

            SignInWithAppleButton(.signIn) { request in
                request.requestedScopes = [.fullName, .email]
            } onCompletion: { result in
                Task {
                    await authViewModel.loginWithApple(result: result)
                }
            }
            .signInWithAppleButtonStyle(colorScheme == .dark ? .white : .black)
            .frame(maxWidth: .infinity)
            .frame(height: 50)
            .clipShape(RoundedRectangle(cornerRadius: 12))

            Button {
                guard let window else { return }
                Task {
                    await authViewModel.login(anchor: window)
                }
            } label: {
                Label("Sign in with SSO", systemImage: "person.badge.key.fill")
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 4)
            }
            .buttonStyle(.borderedProminent)
            .tint(.brown)
            .controlSize(.large)

            if let error = authViewModel.error {
                VStack(spacing: 6) {
                    Text(error.message)
                        .font(.subheadline.weight(.medium))
                        .foregroundStyle(.red)
                        .multilineTextAlignment(.center)

                    if let detail = error.detail {
                        Text(detail)
                            .font(.caption2)
                            .foregroundStyle(.secondary)
                            .multilineTextAlignment(.center)
                    }
                }
                .padding()
                .background(.red.opacity(0.08), in: RoundedRectangle(cornerRadius: 10))
            }
        }
        .padding(32)
        .background(WindowFinder(window: $window))
    }
}

/// Invisible UIView that captures its owning UIWindow for use as an ASWebAuthenticationSession anchor.
private struct WindowFinder: UIViewRepresentable {
    @Binding var window: UIWindow?

    func makeUIView(context: Context) -> UIView {
        let view = UIView()
        DispatchQueue.main.async {
            self.window = view.window
        }
        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        DispatchQueue.main.async {
            self.window = uiView.window
        }
    }
}
