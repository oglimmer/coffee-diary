package com.coffeediary.service;

import com.coffeediary.dto.SieveRequest;
import com.coffeediary.dto.SieveResponse;
import com.coffeediary.entity.Sieve;
import com.coffeediary.entity.User;
import com.coffeediary.repository.SieveRepository;
import com.coffeediary.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SieveServiceTest {

    @Mock
    private SieveRepository sieveRepository;

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private SieveService sieveService;

    private User user;

    @BeforeEach
    void setUp() {
        user = User.builder().id(1L).username("testuser").build();
    }

    @Test
    void findAllByUser_returnsMappedResponses() {
        Sieve s1 = Sieve.builder().id(1L).name("IMS 18g").user(user).build();
        Sieve s2 = Sieve.builder().id(2L).name("VST 20g").user(user).build();
        when(sieveRepository.findAllByUserId(1L)).thenReturn(List.of(s1, s2));

        List<SieveResponse> result = sieveService.findAllByUser(1L);

        assertThat(result).hasSize(2);
        assertThat(result.get(0).getName()).isEqualTo("IMS 18g");
        assertThat(result.get(1).getName()).isEqualTo("VST 20g");
    }

    @Test
    void findAllByUser_emptyList() {
        when(sieveRepository.findAllByUserId(1L)).thenReturn(List.of());

        List<SieveResponse> result = sieveService.findAllByUser(1L);

        assertThat(result).isEmpty();
    }

    @Test
    void create_success() {
        when(userRepository.getReferenceById(1L)).thenReturn(user);
        when(sieveRepository.save(any(Sieve.class))).thenAnswer(invocation -> {
            Sieve sieve = invocation.getArgument(0);
            sieve.setId(10L);
            return sieve;
        });

        SieveResponse response = sieveService.create(1L, new SieveRequest("IMS 18g"));

        assertThat(response.getId()).isEqualTo(10L);
        assertThat(response.getName()).isEqualTo("IMS 18g");
    }

    @Test
    void delete_success() {
        Sieve sieve = Sieve.builder().id(5L).name("IMS 18g").user(user).build();
        when(sieveRepository.findById(5L)).thenReturn(Optional.of(sieve));

        sieveService.delete(1L, 5L);

        verify(sieveRepository).delete(sieve);
    }

    @Test
    void delete_notFound_throwsIllegalArgument() {
        when(sieveRepository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> sieveService.delete(1L, 99L))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessage("Sieve not found");
    }

    @Test
    void delete_wrongUser_throwsSecurityException() {
        User otherUser = User.builder().id(2L).username("other").build();
        Sieve sieve = Sieve.builder().id(5L).name("IMS 18g").user(otherUser).build();
        when(sieveRepository.findById(5L)).thenReturn(Optional.of(sieve));

        assertThatThrownBy(() -> sieveService.delete(1L, 5L))
                .isInstanceOf(SecurityException.class)
                .hasMessage("Not authorized to delete this sieve");

        verify(sieveRepository, never()).delete(any());
    }
}
