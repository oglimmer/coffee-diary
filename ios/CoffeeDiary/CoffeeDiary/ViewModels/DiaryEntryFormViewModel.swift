import Foundation

@Observable
final class DiaryEntryFormViewModel {
    // Form fields
    var date = Date()
    var temperature = 93
    var coffeeId: Int64?
    var sieveId: Int64?
    var grindSize = ""
    var inputWeight = ""
    var outputWeight = ""
    var timeSeconds = ""
    var rating: Int?
    var notes = ""

    // State
    var coffees: [Coffee] = []
    var sieves: [Sieve] = []
    var isSaving = false
    var error: String?

    let editingEntry: DiaryEntry?
    var isEditing: Bool { editingEntry != nil }

    private let api = APIClient.shared

    init(entry: DiaryEntry? = nil) {
        editingEntry = entry
        if let entry {
            date = entry.parsedDate ?? Date()
            temperature = entry.temperature
            coffeeId = entry.coffeeId
            sieveId = entry.sieveId
            grindSize = entry.grindSize.map { String(format: "%.1f", $0) } ?? ""
            inputWeight = entry.inputWeight.map { String(format: "%.1f", $0) } ?? ""
            outputWeight = entry.outputWeight.map { String(format: "%.1f", $0) } ?? ""
            timeSeconds = entry.timeSeconds.map { "\($0)" } ?? ""
            rating = entry.rating
            notes = entry.notes ?? ""
        }
    }

    func loadOptions() async {
        async let loadedCoffees: [Coffee] = api.request("GET", path: "/api/coffees")
        async let loadedSieves: [Sieve] = api.request("GET", path: "/api/sieves")
        coffees = (try? await loadedCoffees) ?? []
        sieves = (try? await loadedSieves) ?? []
    }

    func save() async -> DiaryEntry? {
        isSaving = true
        error = nil

        let request = DiaryEntryRequest(
            dateTime: DiaryEntry.dateFormatter.string(from: date),
            sieveId: sieveId,
            temperature: temperature,
            coffeeId: coffeeId,
            grindSize: Double(grindSize),
            inputWeight: Double(inputWeight),
            outputWeight: Double(outputWeight),
            timeSeconds: Int(timeSeconds),
            rating: rating,
            notes: notes.isEmpty ? nil : notes
        )

        do {
            let entry: DiaryEntry
            if let existing = editingEntry {
                entry = try await api.request("PUT", path: "/api/diary-entries/\(existing.id)", body: request)
            } else {
                entry = try await api.request("POST", path: "/api/diary-entries", body: request)
            }
            isSaving = false
            return entry
        } catch {
            self.error = error.localizedDescription
            isSaving = false
            return nil
        }
    }
}

private extension DiaryEntryRequest {
    init(dateTime: String, sieveId: Int64?, temperature: Int, coffeeId: Int64?,
         grindSize: Double?, inputWeight: Double?, outputWeight: Double?,
         timeSeconds: Int?, rating: Int?, notes: String?) {
        self.dateTime = dateTime
        self.sieveId = sieveId
        self.temperature = temperature
        self.coffeeId = coffeeId
        self.grindSize = grindSize
        self.inputWeight = inputWeight
        self.outputWeight = outputWeight
        self.timeSeconds = timeSeconds
        self.rating = rating
        self.notes = notes
    }
}
