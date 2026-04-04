import Testing
import Foundation
@testable import CoffeeDiary

@Suite("DiaryEntryFilter")
struct DiaryEntryFilterTests {

    @Test("isEmpty returns true when all fields are nil")
    func emptyFilter() {
        let filter = DiaryEntryFilter()
        #expect(filter.isEmpty)
    }

    @Test("isEmpty returns false when coffeeId is set")
    func notEmptyWithCoffeeId() {
        var filter = DiaryEntryFilter()
        filter.coffeeId = 1
        #expect(!filter.isEmpty)
    }

    @Test("isEmpty returns false when sieveId is set")
    func notEmptyWithSieveId() {
        var filter = DiaryEntryFilter()
        filter.sieveId = 5
        #expect(!filter.isEmpty)
    }

    @Test("isEmpty returns false when ratingMin is set")
    func notEmptyWithRating() {
        var filter = DiaryEntryFilter()
        filter.ratingMin = 3
        #expect(!filter.isEmpty)
    }

    @Test("isEmpty returns false when dateFrom is set")
    func notEmptyWithDateFrom() {
        var filter = DiaryEntryFilter()
        filter.dateFrom = Date()
        #expect(!filter.isEmpty)
    }

    @Test("isEmpty returns false when dateTo is set")
    func notEmptyWithDateTo() {
        var filter = DiaryEntryFilter()
        filter.dateTo = Date()
        #expect(!filter.isEmpty)
    }

    @Test("queryItems returns empty array for empty filter")
    func queryItemsEmpty() {
        let filter = DiaryEntryFilter()
        #expect(filter.queryItems().isEmpty)
    }

    @Test("queryItems includes coffeeId")
    func queryItemsCoffeeId() {
        var filter = DiaryEntryFilter()
        filter.coffeeId = 42
        let items = filter.queryItems()
        #expect(items.count == 1)
        #expect(items.first?.name == "coffeeId")
        #expect(items.first?.value == "42")
    }

    @Test("queryItems includes sieveId")
    func queryItemsSieveId() {
        var filter = DiaryEntryFilter()
        filter.sieveId = 7
        let items = filter.queryItems()
        #expect(items.count == 1)
        #expect(items.first?.name == "sieveId")
        #expect(items.first?.value == "7")
    }

    @Test("queryItems includes ratingMin")
    func queryItemsRatingMin() {
        var filter = DiaryEntryFilter()
        filter.ratingMin = 4
        let items = filter.queryItems()
        #expect(items.count == 1)
        #expect(items.first?.name == "ratingMin")
        #expect(items.first?.value == "4")
    }

    @Test("queryItems formats dates using DiaryEntry.dateFormatter")
    func queryItemsDates() {
        let date = DiaryEntry.dateFormatter.date(from: "2025-06-15T10:00:00")!
        var filter = DiaryEntryFilter()
        filter.dateFrom = date
        filter.dateTo = date
        let items = filter.queryItems()
        #expect(items.count == 2)
        #expect(items.first { $0.name == "dateFrom" }?.value == "2025-06-15T10:00:00")
        #expect(items.first { $0.name == "dateTo" }?.value == "2025-06-15T10:00:00")
    }

    @Test("queryItems includes all fields when fully populated")
    func queryItemsAllFields() {
        var filter = DiaryEntryFilter()
        filter.coffeeId = 1
        filter.sieveId = 2
        filter.dateFrom = DiaryEntry.dateFormatter.date(from: "2025-01-01T00:00:00")
        filter.dateTo = DiaryEntry.dateFormatter.date(from: "2025-12-31T23:59:59")
        filter.ratingMin = 3
        let items = filter.queryItems()
        #expect(items.count == 5)
    }
}
