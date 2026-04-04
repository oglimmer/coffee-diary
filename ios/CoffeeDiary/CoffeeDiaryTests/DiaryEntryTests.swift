import Testing
import Foundation
@testable import CoffeeDiary

@Suite("DiaryEntry")
struct DiaryEntryTests {

    @Test("parsedDate returns correct date for valid dateTime")
    func parsedDateValid() {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "2025-03-15T14:30:00",
            sieveId: nil, sieveName: nil, temperature: 93,
            coffeeId: nil, coffeeName: nil, grindSize: nil,
            inputWeight: nil, outputWeight: nil, timeSeconds: nil,
            rating: nil, notes: nil
        )

        let date = entry.parsedDate
        #expect(date != nil)

        let calendar = Calendar(identifier: .gregorian)
        let components = calendar.dateComponents(in: TimeZone(identifier: "UTC")!, from: date!)
        // The formatter has no timezone, so it parses in the current timezone.
        // Just verify year/month/day from the formatter's perspective.
        let formatted = DiaryEntry.dateFormatter.string(from: date!)
        #expect(formatted == "2025-03-15T14:30:00")
    }

    @Test("parsedDate returns nil for invalid dateTime")
    func parsedDateInvalid() {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "not-a-date",
            sieveId: nil, sieveName: nil, temperature: 93,
            coffeeId: nil, coffeeName: nil, grindSize: nil,
            inputWeight: nil, outputWeight: nil, timeSeconds: nil,
            rating: nil, notes: nil
        )
        #expect(entry.parsedDate == nil)
    }

    @Test("dateFormatter round-trips correctly")
    func dateFormatterRoundTrip() {
        let original = "2025-12-31T23:59:59"
        let date = DiaryEntry.dateFormatter.date(from: original)
        #expect(date != nil)
        #expect(DiaryEntry.dateFormatter.string(from: date!) == original)
    }

    @Test("decodes from JSON")
    func decodesFromJSON() throws {
        let json = """
        {
            "id": 42,
            "userId": 7,
            "dateTime": "2025-06-01T09:00:00",
            "temperature": 94,
            "coffeeId": 3,
            "coffeeName": "Ethiopian Yirgacheffe",
            "grindSize": 12.5,
            "inputWeight": 18.0,
            "outputWeight": 36.0,
            "timeSeconds": 28,
            "rating": 4,
            "notes": "Fruity and bright"
        }
        """
        let entry = try JSONDecoder().decode(DiaryEntry.self, from: Data(json.utf8))
        #expect(entry.id == 42)
        #expect(entry.userId == 7)
        #expect(entry.temperature == 94)
        #expect(entry.coffeeId == 3)
        #expect(entry.coffeeName == "Ethiopian Yirgacheffe")
        #expect(entry.grindSize == 12.5)
        #expect(entry.inputWeight == 18.0)
        #expect(entry.outputWeight == 36.0)
        #expect(entry.timeSeconds == 28)
        #expect(entry.rating == 4)
        #expect(entry.notes == "Fruity and bright")
        #expect(entry.sieveId == nil)
        #expect(entry.sieveName == nil)
    }
}
