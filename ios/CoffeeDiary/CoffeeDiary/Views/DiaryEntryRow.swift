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
                if let sieveName = entry.sieveName {
                    Label(sieveName, systemImage: "line.3.horizontal.decrease")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
                Label("\(entry.temperature)\u{00B0}C", systemImage: "thermometer.medium")
                    .font(.caption)
                    .foregroundStyle(.secondary)
                if let time = entry.timeSeconds {
                    Label("\(time)s", systemImage: "timer")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
            }

            if let notes = entry.notes, !notes.isEmpty {
                Text(notes)
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .lineLimit(2)
            }
        }
        .padding(.vertical, 4)
    }

    private var formattedDate: String {
        guard let date = entry.parsedDate else { return entry.dateTime }
        return date.formatted(date: .abbreviated, time: .shortened)
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
