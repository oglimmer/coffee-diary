package com.coffeediary.service;

import com.coffeediary.dto.CoffeeRequest;
import com.coffeediary.dto.CoffeeResponse;
import com.coffeediary.entity.Coffee;
import com.coffeediary.entity.User;
import com.coffeediary.repository.CoffeeRepository;
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
class CoffeeServiceTest {

    @Mock
    private CoffeeRepository coffeeRepository;

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private CoffeeService coffeeService;

    private User user;

    @BeforeEach
    void setUp() {
        user = User.builder().id(1L).username("testuser").build();
    }

    @Test
    void findAllByUser_returnsMappedResponses() {
        Coffee c1 = Coffee.builder().id(1L).name("Ethiopia").user(user).build();
        Coffee c2 = Coffee.builder().id(2L).name("Colombia").user(user).build();
        when(coffeeRepository.findAllByUserId(1L)).thenReturn(List.of(c1, c2));

        List<CoffeeResponse> result = coffeeService.findAllByUser(1L);

        assertThat(result).hasSize(2);
        assertThat(result.get(0).getName()).isEqualTo("Ethiopia");
        assertThat(result.get(1).getName()).isEqualTo("Colombia");
    }

    @Test
    void findAllByUser_emptyList() {
        when(coffeeRepository.findAllByUserId(1L)).thenReturn(List.of());

        List<CoffeeResponse> result = coffeeService.findAllByUser(1L);

        assertThat(result).isEmpty();
    }

    @Test
    void create_success() {
        when(userRepository.getReferenceById(1L)).thenReturn(user);
        when(coffeeRepository.save(any(Coffee.class))).thenAnswer(invocation -> {
            Coffee coffee = invocation.getArgument(0);
            coffee.setId(10L);
            return coffee;
        });

        CoffeeResponse response = coffeeService.create(1L, new CoffeeRequest("Kenya AA"));

        assertThat(response.getId()).isEqualTo(10L);
        assertThat(response.getName()).isEqualTo("Kenya AA");
    }

    @Test
    void delete_success() {
        Coffee coffee = Coffee.builder().id(5L).name("Brazil").user(user).build();
        when(coffeeRepository.findById(5L)).thenReturn(Optional.of(coffee));

        coffeeService.delete(1L, 5L);

        verify(coffeeRepository).delete(coffee);
    }

    @Test
    void delete_notFound_throwsIllegalArgument() {
        when(coffeeRepository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> coffeeService.delete(1L, 99L))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessage("Coffee not found");
    }

    @Test
    void delete_wrongUser_throwsSecurityException() {
        User otherUser = User.builder().id(2L).username("other").build();
        Coffee coffee = Coffee.builder().id(5L).name("Brazil").user(otherUser).build();
        when(coffeeRepository.findById(5L)).thenReturn(Optional.of(coffee));

        assertThatThrownBy(() -> coffeeService.delete(1L, 5L))
                .isInstanceOf(SecurityException.class)
                .hasMessage("Not authorized to delete this coffee");

        verify(coffeeRepository, never()).delete(any());
    }
}
