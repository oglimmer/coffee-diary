import AuthenticationServices
import SwiftUI

struct LoginView: View {
    var authViewModel: AuthViewModel
    @Environment(\.colorScheme) private var colorScheme
    @State private var window: UIWindow?
    @State private var appear = false

    // MARK: - Palette (warm broadsheet)

    private var paper: Color {
        colorScheme == .dark
            ? Color(red: 0.10, green: 0.07, blue: 0.05)
            : Color(red: 0.962, green: 0.922, blue: 0.843)
    }
    private var paperHi: Color {
        colorScheme == .dark
            ? Color(red: 0.14, green: 0.10, blue: 0.07)
            : Color(red: 0.985, green: 0.955, blue: 0.885)
    }
    private var ink: Color {
        colorScheme == .dark
            ? Color(red: 0.965, green: 0.925, blue: 0.845)
            : Color(red: 0.16, green: 0.09, blue: 0.055)
    }
    private var accent: Color {
        Color(red: 0.78, green: 0.33, blue: 0.13) // burnt sienna
    }
    private var muted: Color { ink.opacity(0.55) }
    private var hairline: Color { ink.opacity(0.35) }

    // MARK: - Body

    var body: some View {
        GeometryReader { geo in
            ZStack {
                // Base paper + subtle vignette
                RadialGradient(
                    colors: [paperHi, paper],
                    center: .topLeading,
                    startRadius: 40,
                    endRadius: max(geo.size.width, geo.size.height)
                )
                .ignoresSafeArea()

                GrainOverlay(density: 2600)
                    .opacity(colorScheme == .dark ? 0.22 : 0.18)
                    .blendMode(colorScheme == .dark ? .screen : .multiply)
                    .ignoresSafeArea()
                    .allowsHitTesting(false)

                // Decorative typographic marginalia (rotated italic)
                Text("no. 001  ·  a journal of extraction")
                    .font(.system(size: 11, weight: .regular, design: .serif))
                    .italic()
                    .tracking(1.2)
                    .foregroundStyle(muted)
                    .rotationEffect(.degrees(-90))
                    .fixedSize()
                    .position(x: 18, y: geo.size.height / 2)
                    .opacity(appear ? 1 : 0)

                VStack(spacing: 0) {
                    headerBar
                        .padding(.top, 4)
                        .opacity(appear ? 1 : 0)
                        .offset(y: appear ? 0 : -8)

                    Spacer(minLength: 18)

                    masthead
                        .opacity(appear ? 1 : 0)
                        .offset(y: appear ? 0 : 14)

                    steamComposition
                        .frame(height: 150)
                        .padding(.top, 18)
                        .opacity(appear ? 1 : 0)

                    tagline
                        .padding(.top, 14)
                        .opacity(appear ? 1 : 0)

                    Spacer(minLength: 18)

                    signInBlock
                        .opacity(appear ? 1 : 0)
                        .offset(y: appear ? 0 : 14)

                    if let error = authViewModel.error {
                        errorCard(error)
                            .padding(.top, 14)
                            .transition(.opacity.combined(with: .move(edge: .bottom)))
                    }

                    footer
                        .padding(.top, 18)
                        .opacity(appear ? 1 : 0)
                }
                .padding(.horizontal, 34)
                .padding(.vertical, 22)
            }
        }
        .background(WindowFinder(window: $window))
        .onAppear {
            withAnimation(.easeOut(duration: 1.1).delay(0.05)) {
                appear = true
            }
        }
        .animation(.easeInOut(duration: 0.25), value: authViewModel.error?.message)
    }

    // MARK: - Sections

    private var headerBar: some View {
        HStack(alignment: .firstTextBaseline) {
            VStack(alignment: .leading, spacing: 2) {
                Text("VOL. I")
                    .font(.system(size: 10, weight: .bold, design: .monospaced))
                    .tracking(2)
                Text("№ 001")
                    .font(.system(size: 10, weight: .regular, design: .monospaced))
                    .tracking(1.5)
                    .foregroundStyle(muted)
            }

            Spacer()

            Circle()
                .fill(accent)
                .frame(width: 5, height: 5)

            Spacer()

            VStack(alignment: .trailing, spacing: 2) {
                Text("EST. MMXXV")
                    .font(.system(size: 10, weight: .bold, design: .monospaced))
                    .tracking(2)
                Text("the morning edition")
                    .font(.system(size: 10, weight: .regular, design: .serif))
                    .italic()
                    .foregroundStyle(muted)
            }
        }
        .foregroundStyle(ink)
    }

    private var masthead: some View {
        VStack(spacing: 0) {
            Text("COFFEE")
                .font(.system(size: 68, weight: .black, design: .serif))
                .tracking(-1)
                .foregroundStyle(ink)
                .frame(maxWidth: .infinity, alignment: .leading)

            // Hairline rule with inline label
            HStack(spacing: 10) {
                Rectangle().fill(hairline).frame(height: 0.75)
                Text("— the diary of —")
                    .font(.system(size: 10, weight: .regular, design: .serif))
                    .italic()
                    .tracking(1)
                    .foregroundStyle(muted)
                    .fixedSize()
                Rectangle().fill(hairline).frame(height: 0.75)
            }
            .padding(.vertical, 6)

            HStack(alignment: .lastTextBaseline, spacing: 10) {
                Text("DIARY")
                    .font(.system(size: 68, weight: .black, design: .serif))
                    .tracking(-1)
                    .foregroundStyle(ink)
                Spacer()
                Text("•")
                    .font(.system(size: 18, design: .serif))
                    .foregroundStyle(accent)
                    .baselineOffset(18)
            }
        }
    }

    private var steamComposition: some View {
        ZStack {
            SteamView(color: ink)
                .frame(height: 110)
                .offset(y: -18)

            // Minimal cup glyph built from shapes
            VStack(spacing: 0) {
                Spacer()
                CupShape()
                    .stroke(ink, style: StrokeStyle(lineWidth: 1.6, lineCap: .round, lineJoin: .round))
                    .frame(width: 86, height: 46)
                Rectangle()
                    .fill(ink)
                    .frame(width: 110, height: 1.2)
                    .padding(.top, 3)
            }
        }
    }

    private var tagline: some View {
        Text("Notes on the grind, the pour, and the\nquiet theater of the morning cup.")
            .font(.system(size: 16, weight: .regular, design: .serif))
            .italic()
            .multilineTextAlignment(.center)
            .foregroundStyle(ink.opacity(0.78))
            .lineSpacing(3)
    }

    private var signInBlock: some View {
        VStack(spacing: 14) {
            // Editorial section header
            HStack(spacing: 10) {
                Rectangle().fill(hairline).frame(height: 0.75)
                Text("ENTER THE JOURNAL")
                    .font(.system(size: 10, weight: .bold, design: .monospaced))
                    .tracking(2.5)
                    .foregroundStyle(muted)
                    .fixedSize()
                Rectangle().fill(hairline).frame(height: 0.75)
            }

            SignInWithAppleButton(.signIn) { request in
                request.requestedScopes = [.fullName, .email]
            } onCompletion: { result in
                Task { await authViewModel.loginWithApple(result: result) }
            }
            .signInWithAppleButtonStyle(colorScheme == .dark ? .white : .black)
            .frame(height: 52)
            .clipShape(RoundedRectangle(cornerRadius: 2, style: .continuous))
            .overlay(
                RoundedRectangle(cornerRadius: 2, style: .continuous)
                    .stroke(ink, lineWidth: 1)
            )

            // "or" divider
            HStack(spacing: 12) {
                Rectangle().fill(hairline).frame(height: 0.75)
                Text("or")
                    .font(.system(size: 11, design: .serif))
                    .italic()
                    .foregroundStyle(muted)
                Rectangle().fill(hairline).frame(height: 0.75)
            }
            .padding(.vertical, 2)

            // Custom editorial-styled SSO button
            Button {
                guard let window else { return }
                Task { await authViewModel.login(anchor: window) }
            } label: {
                HStack(spacing: 12) {
                    Text("CONTINUE WITH SSO")
                        .font(.system(size: 13, weight: .bold, design: .monospaced))
                        .tracking(2)
                    Spacer()
                    Image(systemName: "arrow.right")
                        .font(.system(size: 13, weight: .semibold))
                }
                .foregroundStyle(ink)
                .padding(.horizontal, 18)
                .frame(height: 52)
                .frame(maxWidth: .infinity)
                .background(
                    ZStack {
                        Rectangle().fill(paperHi)
                        // Corner ornaments
                        VStack {
                            HStack {
                                CornerTick().stroke(accent, lineWidth: 1.2)
                                    .frame(width: 10, height: 10)
                                Spacer()
                                CornerTick().stroke(accent, lineWidth: 1.2)
                                    .rotationEffect(.degrees(90))
                                    .frame(width: 10, height: 10)
                            }
                            Spacer()
                            HStack {
                                CornerTick().stroke(accent, lineWidth: 1.2)
                                    .rotationEffect(.degrees(270))
                                    .frame(width: 10, height: 10)
                                Spacer()
                                CornerTick().stroke(accent, lineWidth: 1.2)
                                    .rotationEffect(.degrees(180))
                                    .frame(width: 10, height: 10)
                            }
                        }
                        .padding(6)
                    }
                )
                .overlay(
                    Rectangle().stroke(ink, lineWidth: 1)
                )
            }
            .buttonStyle(PressableStyle())
        }
    }

    private func errorCard(_ error: LoginError) -> some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack(spacing: 8) {
                Rectangle().fill(accent).frame(width: 3, height: 12)
                Text("ERRATUM")
                    .font(.system(size: 10, weight: .bold, design: .monospaced))
                    .tracking(2)
                    .foregroundStyle(accent)
            }
            Text(error.message)
                .font(.system(size: 13, weight: .semibold, design: .serif))
                .foregroundStyle(ink)
            if let detail = error.detail {
                Text(detail)
                    .font(.system(size: 11, design: .serif))
                    .italic()
                    .foregroundStyle(muted)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding(14)
        .background(paperHi)
        .overlay(Rectangle().stroke(accent.opacity(0.6), lineWidth: 0.75))
    }

    private var footer: some View {
        VStack(spacing: 6) {
            HStack(spacing: 4) {
                Text("By entering, you accept our")
                    .font(.system(size: 10, design: .serif))
                    .italic()
                    .foregroundStyle(muted)
                Link(destination: URL(string: "https://coffee.oglimmer.com/privacy")!) {
                    Text("Privacy Policy →")
                        .font(.system(size: 10, weight: .semibold, design: .serif))
                        .italic()
                        .underline()
                        .foregroundStyle(ink)
                }
            }
            Text("PRINTED IN SWIFT  ·  ©  OGLIMMER")
                .font(.system(size: 9, weight: .regular, design: .monospaced))
                .tracking(1.8)
                .foregroundStyle(muted)
        }
    }
}

// MARK: - Animated steam (TimelineView + Canvas)

private struct SteamView: View {
    var color: Color

    var body: some View {
        TimelineView(.animation(minimumInterval: 1.0 / 30.0)) { timeline in
            Canvas { ctx, size in
                let t = timeline.date.timeIntervalSinceReferenceDate
                let ribbons = 3
                for i in 0..<ribbons {
                    var path = Path()
                    let baseX = size.width / 2 + CGFloat(i - 1) * 22
                    let startY = size.height
                    path.move(to: CGPoint(x: baseX, y: startY))

                    let steps = 40
                    for s in 1...steps {
                        let p = Double(s) / Double(steps)
                        let y = startY - CGFloat(p) * size.height * 0.98
                        let amp = 10.0 * (0.3 + p) // widens upward
                        let phase = t * 1.15 + Double(i) * 1.3
                        let wave = sin(p * .pi * 3.4 + phase) * amp
                        path.addLine(to: CGPoint(x: baseX + CGFloat(wave), y: y))
                    }

                    // Fade upward by drawing multiple stroked copies with decreasing opacity segments
                    let opacity = 0.38 - Double(i) * 0.09
                    ctx.stroke(
                        path,
                        with: .color(color.opacity(max(opacity, 0.12))),
                        style: StrokeStyle(lineWidth: 1.1, lineCap: .round, lineJoin: .round)
                    )
                }
            }
            .mask(
                LinearGradient(
                    colors: [.clear, .black.opacity(0.2), .black, .black],
                    startPoint: .top,
                    endPoint: .bottom
                )
            )
        }
    }
}

// MARK: - Cup shape

private struct CupShape: Shape {
    func path(in rect: CGRect) -> Path {
        var p = Path()
        let bodyRect = CGRect(x: rect.minX + rect.width * 0.1,
                              y: rect.minY,
                              width: rect.width * 0.72,
                              height: rect.height)

        // Cup body (trapezoid with rounded bottom)
        p.move(to: CGPoint(x: bodyRect.minX, y: bodyRect.minY))
        p.addLine(to: CGPoint(x: bodyRect.minX + bodyRect.width * 0.12,
                              y: bodyRect.maxY - 4))
        p.addQuadCurve(
            to: CGPoint(x: bodyRect.maxX - bodyRect.width * 0.12, y: bodyRect.maxY - 4),
            control: CGPoint(x: bodyRect.midX, y: bodyRect.maxY + 6)
        )
        p.addLine(to: CGPoint(x: bodyRect.maxX, y: bodyRect.minY))

        // Rim line
        p.move(to: CGPoint(x: bodyRect.minX - 3, y: bodyRect.minY))
        p.addLine(to: CGPoint(x: bodyRect.maxX + 3, y: bodyRect.minY))

        // Handle
        let handleStart = CGPoint(x: bodyRect.maxX, y: bodyRect.minY + bodyRect.height * 0.22)
        let handleEnd = CGPoint(x: bodyRect.maxX, y: bodyRect.minY + bodyRect.height * 0.72)
        p.move(to: handleStart)
        p.addCurve(
            to: handleEnd,
            control1: CGPoint(x: bodyRect.maxX + 18, y: bodyRect.minY + bodyRect.height * 0.25),
            control2: CGPoint(x: bodyRect.maxX + 18, y: bodyRect.minY + bodyRect.height * 0.70)
        )

        return p
    }
}

// MARK: - Corner tick (editorial frame ornament)

private struct CornerTick: Shape {
    func path(in rect: CGRect) -> Path {
        var p = Path()
        p.move(to: CGPoint(x: rect.minX, y: rect.maxY))
        p.addLine(to: CGPoint(x: rect.minX, y: rect.minY))
        p.addLine(to: CGPoint(x: rect.maxX, y: rect.minY))
        return p
    }
}

// MARK: - Grain overlay

private struct GrainOverlay: View {
    var density: Int = 2000

    var body: some View {
        Canvas { ctx, size in
            var rng = SystemRandomNumberGenerator()
            for _ in 0..<density {
                let x = Double.random(in: 0...size.width, using: &rng)
                let y = Double.random(in: 0...size.height, using: &rng)
                let o = Double.random(in: 0.0...0.35, using: &rng)
                let s = Double.random(in: 0.4...1.2, using: &rng)
                ctx.fill(
                    Path(CGRect(x: x, y: y, width: s, height: s)),
                    with: .color(.black.opacity(o))
                )
            }
        }
        .allowsHitTesting(false)
    }
}

// MARK: - Pressable button style

private struct PressableStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .scaleEffect(configuration.isPressed ? 0.985 : 1.0)
            .opacity(configuration.isPressed ? 0.85 : 1.0)
            .animation(.easeOut(duration: 0.12), value: configuration.isPressed)
    }
}

// MARK: - Window finder (for OIDC presentation anchor)

/// Invisible UIView that captures its owning UIWindow for use as an ASWebAuthenticationSession anchor.
private struct WindowFinder: UIViewRepresentable {
    @Binding var window: UIWindow?

    func makeUIView(context: Context) -> UIView {
        let view = UIView()
        DispatchQueue.main.async { self.window = view.window }
        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        DispatchQueue.main.async { self.window = uiView.window }
    }
}
