import SwiftUI

struct ManageCoffeesView: View {
    @State private var coffees: [Coffee] = []
    @State private var newName = ""
    @State private var isLoading = false
    @State private var error: String?
    private let api = APIClient.shared

    var body: some View {
        List {
            Section {
                HStack {
                    TextField("New coffee name", text: $newName)
                    Button("Add") {
                        Task { await addItem() }
                    }
                    .disabled(newName.trimmingCharacters(in: .whitespaces).isEmpty)
                }
            }

            Section {
                if coffees.isEmpty && !isLoading {
                    Text("No coffees yet")
                        .foregroundStyle(.secondary)
                } else {
                    ForEach(coffees) { coffee in
                        Text(coffee.name)
                    }
                    .onDelete { indexSet in
                        let items = indexSet.map { coffees[$0] }
                        Task {
                            for item in items {
                                try? await api.requestNoContent("DELETE", path: "/api/coffees/\(item.id)")
                            }
                            await load()
                        }
                    }
                }
            }
        }
        .navigationTitle("Coffees")
        .task { await load() }
        .overlay { if isLoading && coffees.isEmpty { ProgressView() } }
        .alert("Error", isPresented: .init(get: { error != nil }, set: { if !$0 { error = nil } })) {
            Button("OK") {}
        } message: {
            Text(error ?? "")
        }
    }

    private func load() async {
        isLoading = true
        coffees = (try? await api.request("GET", path: "/api/coffees") as [Coffee]) ?? []
        isLoading = false
    }

    private func addItem() async {
        let name = newName.trimmingCharacters(in: .whitespaces)
        guard !name.isEmpty else { return }
        do {
            let _: Coffee = try await api.request("POST", path: "/api/coffees", body: ["name": name])
            newName = ""
            await load()
        } catch {
            self.error = error.localizedDescription
        }
    }
}

struct ManageSievesView: View {
    @State private var sieves: [Sieve] = []
    @State private var newName = ""
    @State private var isLoading = false
    @State private var error: String?
    private let api = APIClient.shared

    var body: some View {
        List {
            Section {
                HStack {
                    TextField("New sieve name", text: $newName)
                    Button("Add") {
                        Task { await addItem() }
                    }
                    .disabled(newName.trimmingCharacters(in: .whitespaces).isEmpty)
                }
            }

            Section {
                if sieves.isEmpty && !isLoading {
                    Text("No sieves yet")
                        .foregroundStyle(.secondary)
                } else {
                    ForEach(sieves) { sieve in
                        Text(sieve.name)
                    }
                    .onDelete { indexSet in
                        let items = indexSet.map { sieves[$0] }
                        Task {
                            for item in items {
                                try? await api.requestNoContent("DELETE", path: "/api/sieves/\(item.id)")
                            }
                            await load()
                        }
                    }
                }
            }
        }
        .navigationTitle("Sieves")
        .task { await load() }
        .overlay { if isLoading && sieves.isEmpty { ProgressView() } }
        .alert("Error", isPresented: .init(get: { error != nil }, set: { if !$0 { error = nil } })) {
            Button("OK") {}
        } message: {
            Text(error ?? "")
        }
    }

    private func load() async {
        isLoading = true
        sieves = (try? await api.request("GET", path: "/api/sieves") as [Sieve]) ?? []
        isLoading = false
    }

    private func addItem() async {
        let name = newName.trimmingCharacters(in: .whitespaces)
        guard !name.isEmpty else { return }
        do {
            let _: Sieve = try await api.request("POST", path: "/api/sieves", body: ["name": name])
            newName = ""
            await load()
        } catch {
            self.error = error.localizedDescription
        }
    }
}
