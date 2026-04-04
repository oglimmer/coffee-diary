import SwiftUI

struct FilterView: View {
    var viewModel: DiaryListViewModel
    @Environment(\.dismiss) private var dismiss

    @State private var selectedCoffeeId: Int64?
    @State private var selectedSieveId: Int64?
    @State private var dateFrom: Date?
    @State private var dateTo: Date?
    @State private var ratingMin: Int?

    @State private var showDateFrom = false
    @State private var showDateTo = false

    var body: some View {
        NavigationStack {
            Form {
                Section("Coffee") {
                    Picker("Coffee", selection: $selectedCoffeeId) {
                        Text("Any").tag(nil as Int64?)
                        ForEach(viewModel.coffees) { coffee in
                            Text(coffee.name).tag(coffee.id as Int64?)
                        }
                    }
                }

                Section("Sieve") {
                    Picker("Sieve", selection: $selectedSieveId) {
                        Text("Any").tag(nil as Int64?)
                        ForEach(viewModel.sieves) { sieve in
                            Text(sieve.name).tag(sieve.id as Int64?)
                        }
                    }
                }

                Section("Date Range") {
                    Toggle("From date", isOn: $showDateFrom)
                    if showDateFrom {
                        DatePicker("From", selection: Binding(
                            get: { dateFrom ?? Date() },
                            set: { dateFrom = $0 }
                        ), displayedComponents: .date)
                    }

                    Toggle("To date", isOn: $showDateTo)
                    if showDateTo {
                        DatePicker("To", selection: Binding(
                            get: { dateTo ?? Date() },
                            set: { dateTo = $0 }
                        ), displayedComponents: .date)
                    }
                }

                Section("Minimum Rating") {
                    Picker("Min rating", selection: $ratingMin) {
                        Text("Any").tag(nil as Int?)
                        ForEach(1...5, id: \.self) { r in
                            HStack(spacing: 2) {
                                ForEach(1...r, id: \.self) { _ in
                                    Image(systemName: "star.fill")
                                        .foregroundStyle(.orange)
                                        .font(.caption2)
                                }
                            }.tag(r as Int?)
                        }
                    }
                }
            }
            .navigationTitle("Filter")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Reset") {
                        Task {
                            await viewModel.clearFilter()
                            dismiss()
                        }
                    }
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("Apply") {
                        let filter = DiaryEntryFilter(
                            coffeeId: selectedCoffeeId,
                            sieveId: selectedSieveId,
                            dateFrom: showDateFrom ? dateFrom : nil,
                            dateTo: showDateTo ? dateTo : nil,
                            ratingMin: ratingMin
                        )
                        Task {
                            await viewModel.applyFilter(filter)
                            dismiss()
                        }
                    }
                }
            }
            .onAppear {
                selectedCoffeeId = viewModel.filter.coffeeId
                selectedSieveId = viewModel.filter.sieveId
                dateFrom = viewModel.filter.dateFrom
                dateTo = viewModel.filter.dateTo
                showDateFrom = viewModel.filter.dateFrom != nil
                showDateTo = viewModel.filter.dateTo != nil
                ratingMin = viewModel.filter.ratingMin
            }
        }
    }
}
