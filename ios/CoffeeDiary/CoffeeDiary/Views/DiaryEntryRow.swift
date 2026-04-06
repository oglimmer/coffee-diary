import SwiftUI

struct DiaryEntryRow: View {
    let entry: DiaryEntry
    @Environment(\.horizontalSizeClass) private var sizeClass

    var body: some View {
        if sizeClass == .regular {
            iPadLayout
        } else {
            compactLayout
        }
    }

    // MARK: - iPhone Layout

    private var compactLayout: some View {
        VStack(alignment: .leading, spacing: 6) {
            HStack {
                Text(formattedDate)
                    .font(.subheadline.bold())
                Spacer()
                if let rating = entry.rating {
                    RatingView(rating: rating)
                }
            }

            if let coffeeName = entry.coffeeName {
                Label(coffeeName, systemImage: "leaf.fill")
                    .font(.subheadline)
                    .foregroundStyle(.brown)
            }

            HStack(spacing: 12) {
                if let grindSize = entry.grindSize {
                    Label(formatNumber(grindSize), systemImage: "gearshape")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
                if entry.inputWeight != nil || entry.outputWeight != nil {
                    HStack(spacing: 4) {
                        Image(systemName: "arrow.right")
                            .font(.caption2)
                            .foregroundStyle(.secondary)
                        Text(weightFlow)
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                if let time = entry.timeSeconds {
                    Label("\(time)s", systemImage: "timer")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
            }
        }
        .padding(.vertical, 4)
    }

    // MARK: - iPad Layout

    private var iPadLayout: some View {
        HStack(alignment: .top, spacing: 0) {
            // Left column: identity
            VStack(alignment: .leading, spacing: 4) {
                Text(formattedDate)
                    .font(.subheadline.bold())

                if let coffeeName = entry.coffeeName {
                    Label(coffeeName, systemImage: "leaf.fill")
                        .font(.subheadline)
                        .foregroundStyle(.brown)
                        .lineLimit(1)
                }

                if let sieveName = entry.sieveName {
                    Label(sieveName, systemImage: "line.3.horizontal.decrease")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                        .lineLimit(1)
                }
            }
            .frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)

            // Center column: brew parameters
            VStack(alignment: .leading, spacing: 4) {
                parameterRow(icon: "thermometer.medium", value: "\(entry.temperature)°C")

                if let grindSize = entry.grindSize {
                    parameterRow(icon: "gearshape", value: formatNumber(grindSize))
                }

                if entry.inputWeight != nil || entry.outputWeight != nil {
                    HStack(spacing: 4) {
                        Image(systemName: "scalemass")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                            .frame(width: 14)
                        Text(weightFlow)
                            .font(.caption)
                            .foregroundStyle(.secondary)
                        if let ratio = brewRatio {
                            Text("(\(ratio))")
                                .font(.caption)
                                .foregroundStyle(.orange)
                        }
                    }
                }

                if let time = entry.timeSeconds {
                    parameterRow(icon: "timer", value: "\(time)s")
                }
            }
            .frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)

            // Right column: rating + notes
            VStack(alignment: .trailing, spacing: 6) {
                if let rating = entry.rating {
                    RatingView(rating: rating)
                }

                if let notes = entry.notes, !notes.isEmpty {
                    Text(notes)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                        .lineLimit(2)
                        .multilineTextAlignment(.trailing)
                }
            }
            .frame(minWidth: 0, maxWidth: .infinity, alignment: .trailing)
        }
        .padding(.vertical, 6)
    }

    private func parameterRow(icon: String, value: String) -> some View {
        HStack(spacing: 4) {
            Image(systemName: icon)
                .font(.caption)
                .foregroundStyle(.secondary)
                .frame(width: 14)
            Text(value)
                .font(.caption)
                .foregroundStyle(.secondary)
        }
    }

    // MARK: - Helpers

    private var formattedDate: String {
        guard let date = entry.parsedDate else { return entry.dateTime }
        return date.formatted(date: .abbreviated, time: .shortened)
    }

    private var weightFlow: String {
        let input = entry.inputWeight.map { formatNumber($0) + "g" } ?? "?"
        let output = entry.outputWeight.map { formatNumber($0) + "g" } ?? "?"
        return "\(input) → \(output)"
    }

    private var brewRatio: String? {
        guard let input = entry.inputWeight, let output = entry.outputWeight, input > 0 else {
            return nil
        }
        return String(format: "1:%.1f", output / input)
    }

    private func formatNumber(_ value: Double) -> String {
        value.truncatingRemainder(dividingBy: 1) == 0
            ? String(format: "%.0f", value)
            : String(format: "%.1f", value)
    }
}

struct RatingView: View {
    let rating: Int

    var body: some View {
        HStack(spacing: 2) {
            ForEach(1...5, id: \.self) { star in
                Image(systemName: star <= rating ? "star.fill" : "star")
                    .font(.caption2)
                    .foregroundStyle(star <= rating ? .orange : .gray.opacity(0.3))
            }
        }
    }
}
