package com.coffeediary;

import com.coffeediary.dto.*;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import jakarta.servlet.http.Cookie;
import org.junit.jupiter.api.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc;
import org.springframework.http.MediaType;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.MvcResult;
import org.testcontainers.containers.MariaDBContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.time.LocalDateTime;

import static org.hamcrest.Matchers.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
@Testcontainers
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class ApiIT {

    @Container
    static MariaDBContainer<?> mariaDB = new MariaDBContainer<>("mariadb:11")
            .withDatabaseName("coffeediary")
            .withUsername("test")
            .withPassword("test");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", mariaDB::getJdbcUrl);
        registry.add("spring.datasource.username", mariaDB::getUsername);
        registry.add("spring.datasource.password", mariaDB::getPassword);
        registry.add("rate-limit.requests-per-second", () -> "1000.0");
    }

    @Autowired
    private MockMvc mockMvc;

    private final ObjectMapper objectMapper = new ObjectMapper()
            .findAndRegisterModules()
            .disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);

    private Cookie getSessionCookie(MvcResult result) {
        for (Cookie cookie : result.getResponse().getCookies()) {
            if ("SESSION".equals(cookie.getName())) {
                return cookie;
            }
        }
        throw new AssertionError("No SESSION cookie in response");
    }

    private Cookie registerAndLogin(String username, String password) throws Exception {
        mockMvc.perform(post("/api/auth/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(new RegisterRequest(username, password))));

        MvcResult loginResult = mockMvc.perform(post("/api/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new LoginRequest(username, password))))
                .andExpect(status().isOk())
                .andReturn();

        return getSessionCookie(loginResult);
    }

    // --- Auth Tests ---

    @Test
    @Order(1)
    void register_success() throws Exception {
        mockMvc.perform(post("/api/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new RegisterRequest("alice", "password123"))))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.username").value("alice"))
                .andExpect(jsonPath("$.id").isNumber());
    }

    @Test
    @Order(2)
    void register_duplicateUsername_returns400() throws Exception {
        mockMvc.perform(post("/api/auth/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(new RegisterRequest("dupuser", "password123"))));

        mockMvc.perform(post("/api/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new RegisterRequest("dupuser", "password123"))))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Username already exists"));
    }

    @Test
    @Order(3)
    void register_validation_shortUsername() throws Exception {
        mockMvc.perform(post("/api/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new RegisterRequest("ab", "password123"))))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error").value("Validation Failed"));
    }

    @Test
    @Order(4)
    void register_validation_shortPassword() throws Exception {
        mockMvc.perform(post("/api/auth/register")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new RegisterRequest("validname", "12345"))))
                .andExpect(status().isBadRequest());
    }

    @Test
    @Order(5)
    void login_success() throws Exception {
        mockMvc.perform(post("/api/auth/register")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(new RegisterRequest("loginuser", "password123"))));

        mockMvc.perform(post("/api/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new LoginRequest("loginuser", "password123"))))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.username").value("loginuser"));
    }

    @Test
    @Order(6)
    void login_badCredentials_returns401() throws Exception {
        mockMvc.perform(post("/api/auth/login")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new LoginRequest("nonexistent", "wrong"))))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @Order(7)
    void me_authenticated() throws Exception {
        Cookie session = registerAndLogin("meuser", "password123");

        mockMvc.perform(get("/api/auth/me").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.username").value("meuser"));
    }

    @Test
    @Order(8)
    void me_unauthenticated_returns401() throws Exception {
        mockMvc.perform(get("/api/auth/me"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @Order(9)
    void logout_success() throws Exception {
        Cookie session = registerAndLogin("logoutuser", "password123");

        mockMvc.perform(post("/api/auth/logout").cookie(session))
                .andExpect(status().isNoContent());

        mockMvc.perform(get("/api/auth/me").cookie(session))
                .andExpect(status().isUnauthorized());
    }

    // --- Coffee CRUD ---

    @Test
    @Order(20)
    void coffee_crud() throws Exception {
        Cookie session = registerAndLogin("coffeeuser", "password123");

        // Create
        MvcResult createResult = mockMvc.perform(post("/api/coffees")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest("Ethiopia Yirgacheffe"))))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.name").value("Ethiopia Yirgacheffe"))
                .andExpect(jsonPath("$.id").isNumber())
                .andReturn();

        Long coffeeId = objectMapper.readTree(createResult.getResponse().getContentAsString()).get("id").asLong();

        // Create another
        mockMvc.perform(post("/api/coffees")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest("Colombia Supremo"))))
                .andExpect(status().isCreated());

        // List
        mockMvc.perform(get("/api/coffees").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(2)));

        // Delete
        mockMvc.perform(delete("/api/coffees/" + coffeeId).cookie(session))
                .andExpect(status().isNoContent());

        // Verify deletion
        mockMvc.perform(get("/api/coffees").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(1)));
    }

    @Test
    @Order(21)
    void coffee_validation_blankName() throws Exception {
        Cookie session = registerAndLogin("coffeevaluser", "password123");

        mockMvc.perform(post("/api/coffees")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest(""))))
                .andExpect(status().isBadRequest());
    }

    @Test
    @Order(22)
    void coffee_unauthenticated_returns401() throws Exception {
        mockMvc.perform(get("/api/coffees"))
                .andExpect(status().isUnauthorized());
    }

    // --- Sieve CRUD ---

    @Test
    @Order(30)
    void sieve_crud() throws Exception {
        Cookie session = registerAndLogin("sieveuser", "password123");

        // Create
        MvcResult createResult = mockMvc.perform(post("/api/sieves")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new SieveRequest("IMS 18g"))))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.name").value("IMS 18g"))
                .andReturn();

        Long sieveId = objectMapper.readTree(createResult.getResponse().getContentAsString()).get("id").asLong();

        // List
        mockMvc.perform(get("/api/sieves").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(1)));

        // Delete
        mockMvc.perform(delete("/api/sieves/" + sieveId).cookie(session))
                .andExpect(status().isNoContent());

        // Verify deletion
        mockMvc.perform(get("/api/sieves").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(0)));
    }

    @Test
    @Order(31)
    void sieve_unauthenticated_returns401() throws Exception {
        mockMvc.perform(get("/api/sieves"))
                .andExpect(status().isUnauthorized());
    }

    // --- Diary Entry CRUD ---

    @Test
    @Order(40)
    void diaryEntry_crud() throws Exception {
        Cookie session = registerAndLogin("diaryuser", "password123");

        // Create coffee and sieve
        MvcResult coffeeResult = mockMvc.perform(post("/api/coffees")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest("Test Coffee"))))
                .andReturn();
        Long coffeeId = objectMapper.readTree(coffeeResult.getResponse().getContentAsString()).get("id").asLong();

        MvcResult sieveResult = mockMvc.perform(post("/api/sieves")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new SieveRequest("Test Sieve"))))
                .andReturn();
        Long sieveId = objectMapper.readTree(sieveResult.getResponse().getContentAsString()).get("id").asLong();

        // Create diary entry
        DiaryEntryRequest entryRequest = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.of(2026, 3, 18, 10, 0))
                .coffeeId(coffeeId)
                .sieveId(sieveId)
                .temperature(94)
                .grindSize(3.5)
                .inputWeight(18.0)
                .outputWeight(36.0)
                .timeSeconds(28)
                .rating(4)
                .notes("Great espresso")
                .build();

        MvcResult createResult = mockMvc.perform(post("/api/diary-entries")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(entryRequest)))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.temperature").value(94))
                .andExpect(jsonPath("$.rating").value(4))
                .andExpect(jsonPath("$.coffeeId").value(coffeeId))
                .andExpect(jsonPath("$.sieveId").value(sieveId))
                .andExpect(jsonPath("$.coffeeName").value("Test Coffee"))
                .andExpect(jsonPath("$.sieveName").value("Test Sieve"))
                .andExpect(jsonPath("$.notes").value("Great espresso"))
                .andReturn();

        Long entryId = objectMapper.readTree(createResult.getResponse().getContentAsString()).get("id").asLong();

        // Get by ID
        mockMvc.perform(get("/api/diary-entries/" + entryId).cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(entryId));

        // Update
        DiaryEntryRequest updateRequest = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.of(2026, 3, 18, 11, 0))
                .temperature(95)
                .rating(5)
                .notes("Even better")
                .build();

        mockMvc.perform(put("/api/diary-entries/" + entryId)
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(updateRequest)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.temperature").value(95))
                .andExpect(jsonPath("$.rating").value(5))
                .andExpect(jsonPath("$.coffeeId").isEmpty())
                .andExpect(jsonPath("$.sieveId").isEmpty());

        // List
        mockMvc.perform(get("/api/diary-entries").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(1)))
                .andExpect(jsonPath("$.totalElements").value(1));

        // Delete
        mockMvc.perform(delete("/api/diary-entries/" + entryId).cookie(session))
                .andExpect(status().isNoContent());

        // Verify deletion
        mockMvc.perform(get("/api/diary-entries").cookie(session))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(0)));
    }

    @Test
    @Order(41)
    void diaryEntry_validation_missingDateTime() throws Exception {
        Cookie session = registerAndLogin("diaryvaluser", "password123");

        DiaryEntryRequest request = DiaryEntryRequest.builder().rating(3).build();

        mockMvc.perform(post("/api/diary-entries")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error").value("Validation Failed"));
    }

    @Test
    @Order(42)
    void diaryEntry_validation_ratingOutOfRange() throws Exception {
        Cookie session = registerAndLogin("ratingvaluser", "password123");

        DiaryEntryRequest request = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.now())
                .rating(6)
                .build();

        mockMvc.perform(post("/api/diary-entries")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    @Test
    @Order(43)
    void diaryEntry_unauthenticated_returns401() throws Exception {
        mockMvc.perform(get("/api/diary-entries"))
                .andExpect(status().isUnauthorized());
    }

    @Test
    @Order(44)
    void diaryEntry_defaultTemperature() throws Exception {
        Cookie session = registerAndLogin("tempdefuser", "password123");

        DiaryEntryRequest request = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.of(2026, 3, 18, 10, 0))
                .build();

        mockMvc.perform(post("/api/diary-entries")
                        .cookie(session)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.temperature").value(93));
    }

    // --- Cross-user data isolation ---

    @Test
    @Order(50)
    void crossUser_coffeeIsolation() throws Exception {
        Cookie sessionA = registerAndLogin("userA", "password123");

        MvcResult coffeeResult = mockMvc.perform(post("/api/coffees")
                        .cookie(sessionA)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest("User A Coffee"))))
                .andReturn();
        Long userACoffeeId = objectMapper.readTree(coffeeResult.getResponse().getContentAsString()).get("id").asLong();

        Cookie sessionB = registerAndLogin("userB", "password123");

        mockMvc.perform(get("/api/coffees").cookie(sessionB))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(0)));

        mockMvc.perform(delete("/api/coffees/" + userACoffeeId).cookie(sessionB))
                .andExpect(status().isForbidden());
    }

    @Test
    @Order(51)
    void crossUser_cannotAccessOtherUsersDiaryEntry() throws Exception {
        Cookie sessionC = registerAndLogin("userC", "password123");

        DiaryEntryRequest request = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.of(2026, 3, 18, 10, 0))
                .temperature(93)
                .rating(4)
                .build();

        MvcResult createResult = mockMvc.perform(post("/api/diary-entries")
                        .cookie(sessionC)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andReturn();
        Long entryId = objectMapper.readTree(createResult.getResponse().getContentAsString()).get("id").asLong();

        Cookie sessionD = registerAndLogin("userD", "password123");

        mockMvc.perform(get("/api/diary-entries/" + entryId).cookie(sessionD))
                .andExpect(status().isForbidden());

        mockMvc.perform(put("/api/diary-entries/" + entryId)
                        .cookie(sessionD)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isForbidden());

        mockMvc.perform(delete("/api/diary-entries/" + entryId).cookie(sessionD))
                .andExpect(status().isForbidden());
    }

    @Test
    @Order(52)
    void crossUser_cannotUseOtherUsersCoffeeInEntry() throws Exception {
        Cookie sessionE = registerAndLogin("userE", "password123");

        MvcResult coffeeResult = mockMvc.perform(post("/api/coffees")
                        .cookie(sessionE)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new CoffeeRequest("UserE Coffee"))))
                .andReturn();
        Long userECoffeeId = objectMapper.readTree(coffeeResult.getResponse().getContentAsString()).get("id").asLong();

        Cookie sessionF = registerAndLogin("userF", "password123");

        DiaryEntryRequest request = DiaryEntryRequest.builder()
                .dateTime(LocalDateTime.of(2026, 3, 18, 10, 0))
                .coffeeId(userECoffeeId)
                .temperature(93)
                .build();

        mockMvc.perform(post("/api/diary-entries")
                        .cookie(sessionF)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    // --- Actuator ---

    @Test
    @Order(60)
    void actuator_healthIsPublic() throws Exception {
        mockMvc.perform(get("/actuator/health"))
                .andExpect(status().isOk());
    }

    @Test
    @Order(61)
    void actuator_prometheusRequiresAuth() throws Exception {
        mockMvc.perform(get("/actuator/prometheus"))
                .andExpect(status().isUnauthorized());
    }
}
