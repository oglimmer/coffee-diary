package com.coffeediary.service;

import com.coffeediary.dto.DiaryEntryRequest;
import com.coffeediary.dto.DiaryEntryResponse;
import com.coffeediary.dto.PageResponse;
import com.coffeediary.entity.Coffee;
import com.coffeediary.entity.DiaryEntry;
import com.coffeediary.entity.Sieve;
import com.coffeediary.entity.User;
import com.coffeediary.repository.CoffeeRepository;
import com.coffeediary.repository.DiaryEntryRepository;
import com.coffeediary.repository.SieveRepository;
import com.coffeediary.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DiaryEntryServiceTest {

    @Mock
    private DiaryEntryRepository diaryEntryRepository;

    @Mock
    private UserRepository userRepository;

    @Mock
    private SieveRepository sieveRepository;

    @Mock
    private CoffeeRepository coffeeRepository;

    @InjectMocks
    private DiaryEntryService diaryEntryService;

    private User user;
    private User otherUser;
    private Coffee coffee;
    private Sieve sieve;
    private LocalDateTime now;

    @BeforeEach
    void setUp() {
        user = User.builder().id(1L).username("testuser").build();
        otherUser = User.builder().id(2L).username("other").build();
        coffee = Coffee.builder().id(10L).name("Ethiopia").user(user).build();
        sieve = Sieve.builder().id(20L).name("IMS 18g").user(user).build();
        now = LocalDateTime.of(2026, 3, 18, 10, 0);
    }

    @Nested
    class FindAll {

        @Test
        @SuppressWarnings("unchecked")
        void returnsPagedResults() {
            DiaryEntry entry = DiaryEntry.builder()
                    .id(1L).user(user).dateTime(now).temperature(93).rating(4).build();
            Page<DiaryEntry> page = new PageImpl<>(List.of(entry), PageRequest.of(0, 20), 1);
            when(diaryEntryRepository.findAll(any(Specification.class), any(Pageable.class))).thenReturn(page);

            PageResponse<DiaryEntryResponse> result = diaryEntryService.findAll(
                    1L, null, null, null, null, null, PageRequest.of(0, 20));

            assertThat(result.getContent()).hasSize(1);
            assertThat(result.getTotalElements()).isEqualTo(1);
            assertThat(result.getTotalPages()).isEqualTo(1);
        }

        @Test
        @SuppressWarnings("unchecked")
        void withFilters() {
            Page<DiaryEntry> page = new PageImpl<>(List.of(), PageRequest.of(0, 20), 0);
            when(diaryEntryRepository.findAll(any(Specification.class), any(Pageable.class))).thenReturn(page);

            PageResponse<DiaryEntryResponse> result = diaryEntryService.findAll(
                    1L, 10L, 20L, now.minusDays(7), now, 3, PageRequest.of(0, 20));

            assertThat(result.getContent()).isEmpty();
            verify(diaryEntryRepository).findAll(any(Specification.class), any(Pageable.class));
        }
    }

    @Nested
    class FindById {

        @Test
        void success() {
            DiaryEntry entry = DiaryEntry.builder()
                    .id(1L).user(user).dateTime(now).temperature(93).coffee(coffee).sieve(sieve).rating(4).build();
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(entry));

            DiaryEntryResponse response = diaryEntryService.findById(1L, 1L);

            assertThat(response.getId()).isEqualTo(1L);
            assertThat(response.getCoffeeName()).isEqualTo("Ethiopia");
            assertThat(response.getSieveName()).isEqualTo("IMS 18g");
            assertThat(response.getTemperature()).isEqualTo(93);
        }

        @Test
        void notFound() {
            when(diaryEntryRepository.findById(99L)).thenReturn(Optional.empty());

            assertThatThrownBy(() -> diaryEntryService.findById(1L, 99L))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Diary entry not found");
        }

        @Test
        void wrongUser_throwsSecurityException() {
            DiaryEntry entry = DiaryEntry.builder().id(1L).user(otherUser).dateTime(now).temperature(93).build();
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(entry));

            assertThatThrownBy(() -> diaryEntryService.findById(1L, 1L))
                    .isInstanceOf(SecurityException.class)
                    .hasMessage("Not authorized to view this entry");
        }
    }

    @Nested
    class Create {

        @Test
        void success_withoutCoffeeAndSieve() {
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(diaryEntryRepository.save(any(DiaryEntry.class))).thenAnswer(invocation -> {
                DiaryEntry e = invocation.getArgument(0);
                e.setId(1L);
                return e;
            });

            DiaryEntryRequest request = DiaryEntryRequest.builder()
                    .dateTime(now).temperature(94).grindSize(3.5).inputWeight(18.0)
                    .outputWeight(36.0).timeSeconds(28).rating(4).notes("Great shot").build();

            DiaryEntryResponse response = diaryEntryService.create(1L, request);

            assertThat(response.getId()).isEqualTo(1L);
            assertThat(response.getTemperature()).isEqualTo(94);
            assertThat(response.getRating()).isEqualTo(4);
            assertThat(response.getNotes()).isEqualTo("Great shot");
            assertThat(response.getCoffeeId()).isNull();
            assertThat(response.getSieveId()).isNull();
        }

        @Test
        void success_withCoffeeAndSieve() {
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(coffeeRepository.findById(10L)).thenReturn(Optional.of(coffee));
            when(sieveRepository.findById(20L)).thenReturn(Optional.of(sieve));
            when(diaryEntryRepository.save(any(DiaryEntry.class))).thenAnswer(invocation -> {
                DiaryEntry e = invocation.getArgument(0);
                e.setId(1L);
                return e;
            });

            DiaryEntryRequest request = DiaryEntryRequest.builder()
                    .dateTime(now).coffeeId(10L).sieveId(20L).temperature(93).rating(5).build();

            DiaryEntryResponse response = diaryEntryService.create(1L, request);

            assertThat(response.getCoffeeId()).isEqualTo(10L);
            assertThat(response.getSieveId()).isEqualTo(20L);
        }

        @Test
        void defaultTemperature_whenNull() {
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(diaryEntryRepository.save(any(DiaryEntry.class))).thenAnswer(invocation -> {
                DiaryEntry e = invocation.getArgument(0);
                e.setId(1L);
                return e;
            });

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).temperature(null).build();

            DiaryEntryResponse response = diaryEntryService.create(1L, request);

            assertThat(response.getTemperature()).isEqualTo(93);
        }

        @Test
        void coffeeNotFound_throwsIllegalArgument() {
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(coffeeRepository.findById(99L)).thenReturn(Optional.empty());

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).coffeeId(99L).build();

            assertThatThrownBy(() -> diaryEntryService.create(1L, request))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Coffee not found");
        }

        @Test
        void sieveNotFound_throwsIllegalArgument() {
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(sieveRepository.findById(99L)).thenReturn(Optional.empty());

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).sieveId(99L).build();

            assertThatThrownBy(() -> diaryEntryService.create(1L, request))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Sieve not found");
        }

        @Test
        void coffeeFromOtherUser_throwsIllegalArgument() {
            Coffee otherCoffee = Coffee.builder().id(30L).name("Other").user(otherUser).build();
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(coffeeRepository.findById(30L)).thenReturn(Optional.of(otherCoffee));

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).coffeeId(30L).build();

            assertThatThrownBy(() -> diaryEntryService.create(1L, request))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Coffee does not belong to user");
        }

        @Test
        void sieveFromOtherUser_throwsIllegalArgument() {
            Sieve otherSieve = Sieve.builder().id(30L).name("Other").user(otherUser).build();
            when(userRepository.getReferenceById(1L)).thenReturn(user);
            when(sieveRepository.findById(30L)).thenReturn(Optional.of(otherSieve));

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).sieveId(30L).build();

            assertThatThrownBy(() -> diaryEntryService.create(1L, request))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Sieve does not belong to user");
        }
    }

    @Nested
    class Update {

        private DiaryEntry existingEntry;

        @BeforeEach
        void setUp() {
            existingEntry = DiaryEntry.builder()
                    .id(1L).user(user).dateTime(now).temperature(93).coffee(coffee).sieve(sieve).build();
        }

        @Test
        void success() {
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(existingEntry));
            when(diaryEntryRepository.save(any(DiaryEntry.class))).thenAnswer(invocation -> invocation.getArgument(0));

            DiaryEntryRequest request = DiaryEntryRequest.builder()
                    .dateTime(now.plusHours(1)).temperature(95).rating(5).notes("Updated").build();

            DiaryEntryResponse response = diaryEntryService.update(1L, 1L, request);

            assertThat(response.getTemperature()).isEqualTo(95);
            assertThat(response.getRating()).isEqualTo(5);
            assertThat(response.getNotes()).isEqualTo("Updated");
            assertThat(response.getCoffeeId()).isNull();
            assertThat(response.getSieveId()).isNull();
        }

        @Test
        void notFound_throwsIllegalArgument() {
            when(diaryEntryRepository.findById(99L)).thenReturn(Optional.empty());

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).build();

            assertThatThrownBy(() -> diaryEntryService.update(1L, 99L, request))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Diary entry not found");
        }

        @Test
        void wrongUser_throwsSecurityException() {
            DiaryEntry otherEntry = DiaryEntry.builder().id(1L).user(otherUser).dateTime(now).temperature(93).build();
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(otherEntry));

            DiaryEntryRequest request = DiaryEntryRequest.builder().dateTime(now).build();

            assertThatThrownBy(() -> diaryEntryService.update(1L, 1L, request))
                    .isInstanceOf(SecurityException.class)
                    .hasMessage("Not authorized to update this entry");
        }

        @Test
        void clearsSieveAndCoffeeWhenNull() {
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(existingEntry));
            when(diaryEntryRepository.save(any(DiaryEntry.class))).thenAnswer(invocation -> invocation.getArgument(0));

            DiaryEntryRequest request = DiaryEntryRequest.builder()
                    .dateTime(now).temperature(93).build();

            DiaryEntryResponse response = diaryEntryService.update(1L, 1L, request);

            assertThat(response.getCoffeeId()).isNull();
            assertThat(response.getSieveId()).isNull();
        }
    }

    @Nested
    class Delete {

        @Test
        void success() {
            DiaryEntry entry = DiaryEntry.builder().id(1L).user(user).dateTime(now).temperature(93).build();
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(entry));

            diaryEntryService.delete(1L, 1L);

            verify(diaryEntryRepository).delete(entry);
        }

        @Test
        void notFound_throwsIllegalArgument() {
            when(diaryEntryRepository.findById(99L)).thenReturn(Optional.empty());

            assertThatThrownBy(() -> diaryEntryService.delete(1L, 99L))
                    .isInstanceOf(IllegalArgumentException.class)
                    .hasMessage("Diary entry not found");
        }

        @Test
        void wrongUser_throwsSecurityException() {
            DiaryEntry entry = DiaryEntry.builder().id(1L).user(otherUser).dateTime(now).temperature(93).build();
            when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(entry));

            assertThatThrownBy(() -> diaryEntryService.delete(1L, 1L))
                    .isInstanceOf(SecurityException.class)
                    .hasMessage("Not authorized to delete this entry");

            verify(diaryEntryRepository, never()).delete(any(DiaryEntry.class));
        }
    }

    @Test
    void toResponse_nullCoffeeAndSieve() {
        DiaryEntry entry = DiaryEntry.builder()
                .id(1L).user(user).dateTime(now).temperature(93).build();
        when(diaryEntryRepository.findById(1L)).thenReturn(Optional.of(entry));

        DiaryEntryResponse response = diaryEntryService.findById(1L, 1L);

        assertThat(response.getCoffeeId()).isNull();
        assertThat(response.getCoffeeName()).isNull();
        assertThat(response.getSieveId()).isNull();
        assertThat(response.getSieveName()).isNull();
    }
}
