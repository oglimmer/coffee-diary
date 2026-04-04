import Testing
import Foundation
@testable import CoffeeDiary

@Suite("PagedResponse")
struct PagedResponseTests {

    @Test("decodes paged response of coffees")
    func decodesCoffees() throws {
        let json = """
        {
            "content": [
                {"id": 1, "name": "Ethiopian Sidamo"},
                {"id": 2, "name": "Colombian Supremo"}
            ],
            "totalElements": 15,
            "totalPages": 8,
            "number": 0,
            "size": 2
        }
        """
        let response = try JSONDecoder().decode(PagedResponse<Coffee>.self, from: Data(json.utf8))
        #expect(response.content.count == 2)
        #expect(response.content[0].name == "Ethiopian Sidamo")
        #expect(response.content[1].id == 2)
        #expect(response.totalElements == 15)
        #expect(response.totalPages == 8)
        #expect(response.number == 0)
        #expect(response.size == 2)
    }

    @Test("decodes empty page")
    func decodesEmptyPage() throws {
        let json = """
        {
            "content": [],
            "totalElements": 0,
            "totalPages": 0,
            "number": 0,
            "size": 20
        }
        """
        let response = try JSONDecoder().decode(PagedResponse<Sieve>.self, from: Data(json.utf8))
        #expect(response.content.isEmpty)
        #expect(response.totalElements == 0)
    }
}
