import Foundation

struct DiaryEntry: Codable, Identifiable, Hashable {
    let id: Int64
    let userId: Int64
    let dateTime: String
    let sieveId: Int64?
    let sieveName: String?
    let temperature: Int
    let coffeeId: Int64?
    let coffeeName: String?
    let grindSize: Double?
    let inputWeight: Double?
    let outputWeight: Double?
    let timeSeconds: Int?
    let rating: Int?
    let notes: String?

    var parsedDate: Date? {
        DiaryEntry.dateFormatter.date(from: dateTime)
    }

    static let dateFormatter: DateFormatter = {
        let f = DateFormatter()
        f.dateFormat = "yyyy-MM-dd'T'HH:mm:ss"
        f.locale = Locale(identifier: "en_US_POSIX")
        return f
    }()
}
