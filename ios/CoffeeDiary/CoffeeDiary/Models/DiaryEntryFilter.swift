import Foundation

struct DiaryEntryFilter {
    var coffeeId: Int64?
    var sieveId: Int64?
    var dateFrom: Date?
    var dateTo: Date?
    var ratingMin: Int?

    var isEmpty: Bool {
        coffeeId == nil && sieveId == nil && dateFrom == nil && dateTo == nil && ratingMin == nil
    }

    func queryItems() -> [URLQueryItem] {
        var items: [URLQueryItem] = []
        if let coffeeId { items.append(.init(name: "coffeeId", value: "\(coffeeId)")) }
        if let sieveId { items.append(.init(name: "sieveId", value: "\(sieveId)")) }
        if let dateFrom {
            items.append(.init(name: "dateFrom", value: DiaryEntry.dateFormatter.string(from: dateFrom)))
        }
        if let dateTo {
            items.append(.init(name: "dateTo", value: DiaryEntry.dateFormatter.string(from: dateTo)))
        }
        if let ratingMin { items.append(.init(name: "ratingMin", value: "\(ratingMin)")) }
        return items
    }
}
