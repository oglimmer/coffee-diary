import Foundation

struct DiaryEntryRequest: Encodable {
    var dateTime: String
    var sieveId: Int64?
    var temperature: Int
    var coffeeId: Int64?
    var grindSize: Double?
    var inputWeight: Double?
    var outputWeight: Double?
    var timeSeconds: Int?
    var rating: Int?
    var notes: String?

    init(from entry: DiaryEntry? = nil) {
        if let entry {
            dateTime = entry.dateTime
            sieveId = entry.sieveId
            temperature = entry.temperature
            coffeeId = entry.coffeeId
            grindSize = entry.grindSize
            inputWeight = entry.inputWeight
            outputWeight = entry.outputWeight
            timeSeconds = entry.timeSeconds
            rating = entry.rating
            notes = entry.notes
        } else {
            dateTime = DiaryEntry.dateFormatter.string(from: Date())
            temperature = 93
        }
    }
}
