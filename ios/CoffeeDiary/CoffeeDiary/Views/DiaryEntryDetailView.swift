import SwiftUI

struct DiaryEntryDetailView: View {
    let entry: DiaryEntry

    var body: some View {
        List {
            Section("Brew") {
                row("Date", value: formattedDate)
                row("Temperature", value: "\(entry.temperature)\u{00B0}C")
                if let coffeeName = entry.coffeeName {
                    row("Coffee", value: coffeeName)
                }
                if let sieveName = entry.sieveName {
                    row("Sieve", value: sieveName)
                }
                if let grindSize = entry.grindSize {
                    row("Grind Size", value: String(format: "%.1f", grindSize))
                }
            }

            Section("Weight & Time") {
                if let input = entry.inputWeight {
                    row("Input", value: String(format: "%.1f g", input))
                }
                if let output = entry.outputWeight {
                    row("Output", value: String(format: "%.1f g", output))
                }
                if let input = entry.inputWeight, let output = entry.outputWeight, input > 0 {
                    row("Ratio", value: String(format: "1:%.1f", output / input))
                }
                if let time = entry.timeSeconds {
                    row("Time", value: "\(time)s")
                }
            }

            if entry.rating != nil || entry.notes != nil {
                Section("Tasting") {
                    if let rating = entry.rating {
                        HStack {
                            Text("Rating")
                                .foregroundStyle(.secondary)
                            Spacer()
                            RatingView(rating: rating)
                        }
                    }
                    if let notes = entry.notes, !notes.isEmpty {
                        VStack(alignment: .leading, spacing: 4) {
                            Text("Notes")
                                .font(.caption)
                                .foregroundStyle(.secondary)
                            Text(notes)
                        }
                    }
                }
            }
        }
        .navigationTitle("Brew Details")
        .navigationBarTitleDisplayMode(.inline)
    }

    private func row(_ label: String, value: String) -> some View {
        HStack {
            Text(label)
                .foregroundStyle(.secondary)
            Spacer()
            Text(value)
        }
    }

    private var formattedDate: String {
        guard let date = entry.parsedDate else { return entry.dateTime }
        return date.formatted(date: .abbreviated, time: .shortened)
    }
}
