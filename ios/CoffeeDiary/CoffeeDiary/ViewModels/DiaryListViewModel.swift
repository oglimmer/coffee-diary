import Foundation

@Observable
final class DiaryListViewModel {
    var entries: [DiaryEntry] = []
    var isLoading = false
    var error: String?

    var currentPage = 0
    var totalPages = 0
    var hasMore: Bool { currentPage + 1 < totalPages }

    var filter = DiaryEntryFilter()
    var coffees: [Coffee] = []
    var sieves: [Sieve] = []

    private let api = APIClient.shared

    func loadEntries(reset: Bool = false) async {
        if reset { currentPage = 0 }
        isLoading = true
        error = nil

        do {
            var queryItems = [
                URLQueryItem(name: "page", value: "\(currentPage)"),
                URLQueryItem(name: "size", value: "20"),
                URLQueryItem(name: "sort", value: "dateTime,desc"),
            ]
            queryItems.append(contentsOf: filter.queryItems())

            let response: PagedResponse<DiaryEntry> = try await api.request("GET", path: "/api/diary-entries", queryItems: queryItems)

            if reset {
                entries = response.content
            } else {
                entries.append(contentsOf: response.content)
            }
            totalPages = response.totalPages
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func loadNextPage() async {
        guard hasMore, !isLoading else { return }
        currentPage += 1
        await loadEntries()
    }

    func refresh() async {
        await loadEntries(reset: true)
    }

    func deleteEntry(_ entry: DiaryEntry) async {
        do {
            try await api.requestNoContent("DELETE", path: "/api/diary-entries/\(entry.id)")
            entries.removeAll { $0.id == entry.id }
        } catch {
            self.error = error.localizedDescription
        }
    }

    func loadFilterOptions() async {
        async let loadedCoffees: [Coffee] = api.request("GET", path: "/api/coffees")
        async let loadedSieves: [Sieve] = api.request("GET", path: "/api/sieves")
        coffees = (try? await loadedCoffees) ?? []
        sieves = (try? await loadedSieves) ?? []
    }

    func applyFilter(_ newFilter: DiaryEntryFilter) async {
        filter = newFilter
        await loadEntries(reset: true)
    }

    func clearFilter() async {
        filter = DiaryEntryFilter()
        await loadEntries(reset: true)
    }
}
