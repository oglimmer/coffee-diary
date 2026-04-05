import SwiftUI

struct DiaryListView: View {
    var authViewModel: AuthViewModel
    @State private var viewModel = DiaryListViewModel()
    @State private var showFilter = false
    @State private var showNewEntry = false
    @State private var editingEntry: DiaryEntry?
    @State private var showSettings = false
    @State private var showDeleteAccountConfirm = false

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.entries.isEmpty && viewModel.isLoading {
                    ProgressView("Loading entries...")
                } else if viewModel.entries.isEmpty && !viewModel.isLoading {
                    ContentUnavailableView(
                        "No Entries",
                        systemImage: "cup.and.saucer",
                        description: Text(viewModel.filter.isEmpty ? "Start tracking your brews!" : "No entries match your filter.")
                    )
                } else {
                    entryList
                }
            }
            .navigationTitle("Coffee Diary")
            .toolbar {
                ToolbarItem(placement: .topBarLeading) {
                    Menu {
                        NavigationLink(destination: ManageCoffeesView()) {
                            Label("Coffees", systemImage: "leaf")
                        }
                        NavigationLink(destination: ManageSievesView()) {
                            Label("Sieves", systemImage: "line.3.horizontal.decrease")
                        }
                        Divider()
                        Button(role: .destructive) {
                            Task { await authViewModel.logout() }
                        } label: {
                            Label("Logout", systemImage: "rectangle.portrait.and.arrow.right")
                        }
                        Button(role: .destructive) {
                            showDeleteAccountConfirm = true
                        } label: {
                            Label("Delete Account", systemImage: "person.crop.circle.badge.xmark")
                        }
                    } label: {
                        Image(systemName: "ellipsis.circle")
                    }
                }
                ToolbarItem(placement: .topBarTrailing) {
                    HStack(spacing: 12) {
                        Button {
                            showFilter = true
                        } label: {
                            Image(systemName: viewModel.filter.isEmpty ? "line.3.horizontal.decrease.circle" : "line.3.horizontal.decrease.circle.fill")
                        }
                        Button {
                            showNewEntry = true
                        } label: {
                            Image(systemName: "plus")
                        }
                    }
                }
            }
            .sheet(isPresented: $showFilter) {
                FilterView(viewModel: viewModel)
            }
            .sheet(isPresented: $showNewEntry) {
                DiaryEntryFormView { _ in
                    Task { await viewModel.refresh() }
                }
            }
            .sheet(item: $editingEntry) { entry in
                DiaryEntryFormView(entry: entry) { _ in
                    Task { await viewModel.refresh() }
                }
            }
            .task {
                await viewModel.loadFilterOptions()
                await viewModel.refresh()
            }
            .refreshable {
                await viewModel.refresh()
            }
            .alert("Error", isPresented: .init(get: { viewModel.error != nil }, set: { if !$0 { viewModel.error = nil } })) {
                Button("OK") { viewModel.error = nil }
            } message: {
                Text(viewModel.error ?? "")
            }
            .alert("Delete Account?", isPresented: $showDeleteAccountConfirm) {
                Button("Cancel", role: .cancel) {}
                Button("Delete Account", role: .destructive) {
                    Task { await authViewModel.deleteAccount() }
                }
            } message: {
                Text("This will permanently delete your account, all your diary entries, coffees, and sieves. This action cannot be undone.")
            }
            .alert(
                "Account Deletion Failed",
                isPresented: .init(
                    get: { authViewModel.error != nil },
                    set: { if !$0 { authViewModel.error = nil } }
                )
            ) {
                Button("OK") { authViewModel.error = nil }
            } message: {
                Text(authViewModel.error?.detail ?? authViewModel.error?.message ?? "")
            }
        }
    }

    private var entryList: some View {
        List {
            ForEach(viewModel.entries) { entry in
                NavigationLink(destination: DiaryEntryDetailView(entry: entry)) {
                    DiaryEntryRow(entry: entry)
                }
                .swipeActions(edge: .trailing) {
                    Button(role: .destructive) {
                        Task { await viewModel.deleteEntry(entry) }
                    } label: {
                        Label("Delete", systemImage: "trash")
                    }
                    Button {
                        editingEntry = entry
                    } label: {
                        Label("Edit", systemImage: "pencil")
                    }
                    .tint(.orange)
                }
                .task {
                    if entry.id == viewModel.entries.last?.id {
                        await viewModel.loadNextPage()
                    }
                }
            }

            if viewModel.isLoading && !viewModel.entries.isEmpty {
                HStack {
                    Spacer()
                    ProgressView()
                    Spacer()
                }
            }
        }
        .listStyle(.plain)
    }
}
