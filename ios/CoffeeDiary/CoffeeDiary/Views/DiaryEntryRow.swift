import SwiftUI

struct DiaryEntryRow: View {
    let entry: DiaryEntry

    var body: some View {
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

    private var formattedDate: String {
        guard let date = entry.parsedDate else { return entry.dateTime }
        return date.formatted(date: .abbreviated, time: .shortened)
    }

    private var weightFlow: String {
        let input = entry.inputWeight.map { formatNumber($0) + "g" } ?? "?"
        let output = entry.outputWeight.map { formatNumber($0) + "g" } ?? "?"
        return "\(input) → \(output)"
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
