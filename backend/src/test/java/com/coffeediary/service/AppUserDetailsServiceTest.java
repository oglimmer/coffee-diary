package com.coffeediary.service;

import com.coffeediary.entity.User;
import com.coffeediary.repository.UserRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class AppUserDetailsServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private AppUserDetailsService appUserDetailsService;

    @Test
    void loadUserByUsername_success() {
        User user = User.builder().id(1L).username("testuser").password("encoded").build();
        when(userRepository.findByUsername("testuser")).thenReturn(Optional.of(user));

        UserDetails details = appUserDetailsService.loadUserByUsername("testuser");

        assertThat(details).isInstanceOf(AppUserDetails.class);
        assertThat(details.getUsername()).isEqualTo("testuser");
        assertThat(details.getPassword()).isEqualTo("encoded");
        assertThat(((AppUserDetails) details).getId()).isEqualTo(1L);
        assertThat(details.getAuthorities()).hasSize(1);
        assertThat(details.getAuthorities().iterator().next().getAuthority()).isEqualTo("ROLE_USER");
    }

    @Test
    void loadUserByUsername_notFound_throwsUsernameNotFoundException() {
        when(userRepository.findByUsername("unknown")).thenReturn(Optional.empty());

        assertThatThrownBy(() -> appUserDetailsService.loadUserByUsername("unknown"))
                .isInstanceOf(UsernameNotFoundException.class)
                .hasMessage("User not found: unknown");
    }
}
