import SwiftUI

struct DiaryEntryFormView: View {
    @State var viewModel: DiaryEntryFormViewModel
    @Environment(\.dismiss) private var dismiss
    var onSaved: (DiaryEntry) -> Void

    init(entry: DiaryEntry? = nil, onSaved: @escaping (DiaryEntry) -> Void) {
        _viewModel = State(initialValue: DiaryEntryFormViewModel(entry: entry))
        self.onSaved = onSaved
    }

    var body: some View {
        NavigationStack {
            Form {
                Section("Date & Temperature") {
                    DatePicker("Date", selection: $viewModel.date)
                    Stepper("Temperature: \(viewModel.temperature)\u{00B0}C", value: $viewModel.temperature, in: 70...100)
                }

                Section("Coffee & Sieve") {
                    Picker("Coffee", selection: $viewModel.coffeeId) {
                        Text("None").tag(nil as Int64?)
                        ForEach(viewModel.coffees) { coffee in
                            Text(coffee.name).tag(coffee.id as Int64?)
                        }
                    }
                    Picker("Sieve", selection: $viewModel.sieveId) {
                        Text("None").tag(nil as Int64?)
                        ForEach(viewModel.sieves) { sieve in
                            Text(sieve.name).tag(sieve.id as Int64?)
                        }
                    }
                }

                Section("Grind & Weight") {
                    HStack {
                        Text("Grind Size")
                        Spacer()
                        TextField("e.g. 8.5", text: $viewModel.grindSize)
                            .keyboardType(.decimalPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }
                    HStack {
                        Text("Input (g)")
                        Spacer()
                        TextField("e.g. 18.0", text: $viewModel.inputWeight)
                            .keyboardType(.decimalPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }
                    HStack {
                        Text("Output (g)")
                        Spacer()
                        TextField("e.g. 36.0", text: $viewModel.outputWeight)
                            .keyboardType(.decimalPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }
                    HStack {
                        Text("Time (seconds)")
                        Spacer()
                        TextField("e.g. 30", text: $viewModel.timeSeconds)
                            .keyboardType(.numberPad)
                            .multilineTextAlignment(.trailing)
                            .frame(width: 80)
                    }
                }

                Section("Tasting") {
                    RatingPicker(rating: $viewModel.rating)
                    TextField("Notes", text: $viewModel.notes, axis: .vertical)
                        .lineLimit(3...6)
                }
            }
            .navigationTitle(viewModel.isEditing ? "Edit Entry" : "New Entry")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("Cancel") { dismiss() }
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("Save") {
                        Task {
                            if let entry = await viewModel.save() {
                                onSaved(entry)
                                dismiss()
                            }
                        }
                    }
                    .disabled(viewModel.isSaving)
                }
            }
            .task {
                await viewModel.loadOptions()
            }
            .alert("Error", isPresented: .init(get: { viewModel.error != nil }, set: { if !$0 { viewModel.error = nil } })) {
                Button("OK") { viewModel.error = nil }
            } message: {
                Text(viewModel.error ?? "")
            }
        }
    }
}

struct RatingPicker: View {
    @Binding var rating: Int?

    var body: some View {
        HStack {
            Text("Rating")
            Spacer()
            HStack(spacing: 8) {
                ForEach(1...5, id: \.self) { star in
                    Button {
                        rating = rating == star ? nil : star
                    } label: {
                        Image(systemName: star <= (rating ?? 0) ? "star.fill" : "star")
                            .foregroundStyle(star <= (rating ?? 0) ? .orange : .gray.opacity(0.4))
                            .font(.title3)
                    }
                    .buttonStyle(.plain)
                }
            }
        }
    }
}
