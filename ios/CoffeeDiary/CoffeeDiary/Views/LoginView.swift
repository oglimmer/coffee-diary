import SwiftUI

struct LoginView: View {
    var authViewModel: AuthViewModel
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
                Text(error)
                    .font(.caption)
                    .foregroundStyle(.red)
                    .multilineTextAlignment(.center)
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
