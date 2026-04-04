import Testing
import Foundation
@testable import CoffeeDiary

@Suite("DiaryEntryFormViewModel")
struct DiaryEntryFormViewModelTests {

    @Test("init without entry sets defaults")
    func initDefaults() {
        let vm = DiaryEntryFormViewModel()
        #expect(vm.temperature == 93)
        #expect(vm.coffeeId == nil)
        #expect(vm.sieveId == nil)
        #expect(vm.grindSize == "")
        #expect(vm.inputWeight == "")
        #expect(vm.outputWeight == "")
        #expect(vm.timeSeconds == "")
        #expect(vm.rating == nil)
        #expect(vm.notes == "")
        #expect(vm.editingEntry == nil)
        #expect(!vm.isEditing)
        #expect(!vm.isSaving)
        #expect(vm.error == nil)
    }

    @Test("init with entry populates all fields")
    func initFromEntry() {
        let entry = DiaryEntry(
            id: 10, userId: 1, dateTime: "2025-06-15T08:30:00",
            sieveId: 2, sieveName: "IMS", temperature: 96,
            coffeeId: 5, coffeeName: "Kenya AA", grindSize: 14.5,
            inputWeight: 18.5, outputWeight: 37.0, timeSeconds: 30,
            rating: 5, notes: "Perfect"
        )

        let vm = DiaryEntryFormViewModel(entry: entry)
        #expect(vm.temperature == 96)
        #expect(vm.coffeeId == 5)
        #expect(vm.sieveId == 2)
        #expect(vm.grindSize == "14.5")
        #expect(vm.inputWeight == "18.5")
        #expect(vm.outputWeight == "37.0")
        #expect(vm.timeSeconds == "30")
        #expect(vm.rating == 5)
        #expect(vm.notes == "Perfect")
        #expect(vm.isEditing)
        #expect(vm.editingEntry?.id == 10)
    }

    @Test("init with entry having nil optionals uses empty strings")
    func initFromEntryNilOptionals() {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "2025-01-01T12:00:00",
            sieveId: nil, sieveName: nil, temperature: 93,
            coffeeId: nil, coffeeName: nil, grindSize: nil,
            inputWeight: nil, outputWeight: nil, timeSeconds: nil,
            rating: nil, notes: nil
        )

        let vm = DiaryEntryFormViewModel(entry: entry)
        #expect(vm.grindSize == "")
        #expect(vm.inputWeight == "")
        #expect(vm.outputWeight == "")
        #expect(vm.timeSeconds == "")
        #expect(vm.notes == "")
        #expect(vm.rating == nil)
        #expect(vm.coffeeId == nil)
        #expect(vm.sieveId == nil)
    }

    @Test("parsedDate is set from entry dateTime")
    func dateFromEntry() {
        let entry = DiaryEntry(
            id: 1, userId: 1, dateTime: "2025-06-15T08:30:00",
            sieveId: nil, sieveName: nil, temperature: 93,
            coffeeId: nil, coffeeName: nil, grindSize: nil,
            inputWeight: nil, outputWeight: nil, timeSeconds: nil,
            rating: nil, notes: nil
        )

        let vm = DiaryEntryFormViewModel(entry: entry)
        let formatted = DiaryEntry.dateFormatter.string(from: vm.date)
        #expect(formatted == "2025-06-15T08:30:00")
    }
}
