import Testing
import Foundation
@testable import CoffeeDiary

@Suite("DiaryEntryRequest")
struct DiaryEntryRequestTests {

    @Test("init without entry sets defaults")
    func initDefaults() {
        let request = DiaryEntryRequest()
        #expect(request.temperature == 93)
        #expect(request.sieveId == nil)
        #expect(request.coffeeId == nil)
        #expect(request.grindSize == nil)
        #expect(request.inputWeight == nil)
        #expect(request.outputWeight == nil)
        #expect(request.timeSeconds == nil)
        #expect(request.rating == nil)
        #expect(request.notes == nil)
        // dateTime should be a valid date string
        #expect(DiaryEntry.dateFormatter.date(from: request.dateTime) != nil)
    }

    @Test("init from entry copies all fields")
    func initFromEntry() {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "2025-03-15T14:30:00",
            sieveId: 5, sieveName: "VST 18g", temperature: 95,
            coffeeId: 3, coffeeName: "Ethiopia", grindSize: 12.5,
            inputWeight: 18.0, outputWeight: 36.0, timeSeconds: 28,
            rating: 4, notes: "Great shot"
        )

        let request = DiaryEntryRequest(from: entry)
        #expect(request.dateTime == "2025-03-15T14:30:00")
        #expect(request.sieveId == 5)
        #expect(request.temperature == 95)
        #expect(request.coffeeId == 3)
        #expect(request.grindSize == 12.5)
        #expect(request.inputWeight == 18.0)
        #expect(request.outputWeight == 36.0)
        #expect(request.timeSeconds == 28)
        #expect(request.rating == 4)
        #expect(request.notes == "Great shot")
    }

    @Test("encodes to JSON with correct keys")
    func encodesToJSON() throws {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "2025-03-15T14:30:00",
            sieveId: nil, sieveName: nil, temperature: 93,
            coffeeId: nil, coffeeName: nil, grindSize: nil,
            inputWeight: nil, outputWeight: nil, timeSeconds: nil,
            rating: 3, notes: nil
        )
        let request = DiaryEntryRequest(from: entry)
        let data = try JSONEncoder().encode(request)
        let dict = try JSONDecoder().decode([String: AnyCodable].self, from: data)
        #expect(dict["dateTime"]?.stringValue == "2025-03-15T14:30:00")
        #expect(dict["temperature"]?.intValue == 93)
        #expect(dict["rating"]?.intValue == 3)
    }
}

/// Minimal type-erased Codable for verifying JSON structure in tests.
private struct AnyCodable: Decodable {
    let value: Any

    var stringValue: String? { value as? String }
    var intValue: Int? { value as? Int }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let v = try? container.decode(Int.self) { value = v }
        else if let v = try? container.decode(Double.self) { value = v }
        else if let v = try? container.decode(String.self) { value = v }
        else if let v = try? container.decode(Bool.self) { value = v }
        else if container.decodeNil() { value = NSNull() }
        else { throw DecodingError.dataCorruptedError(in: container, debugDescription: "Unsupported type") }
    }
}
