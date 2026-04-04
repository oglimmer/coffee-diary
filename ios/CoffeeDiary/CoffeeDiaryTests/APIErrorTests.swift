import Testing
import Foundation
@testable import CoffeeDiary

@Suite("APIError")
struct APIErrorTests {

    @Test("unauthorized description")
    func unauthorized() {
        let error = APIError.unauthorized
        #expect(error.errorDescription == "Authentication required")
    }

    @Test("badRequest includes message")
    func badRequest() {
        let error = APIError.badRequest("Invalid coffee ID")
        #expect(error.errorDescription == "Invalid coffee ID")
    }

    @Test("serverError includes status code")
    func serverError() {
        let error = APIError.serverError(500)
        #expect(error.errorDescription == "Server error (500)")
    }

    @Test("networkError wraps underlying error description")
    func networkError() {
        let underlying = URLError(.notConnectedToInternet)
        let error = APIError.networkError(underlying)
        #expect(error.errorDescription == underlying.localizedDescription)
    }
}
