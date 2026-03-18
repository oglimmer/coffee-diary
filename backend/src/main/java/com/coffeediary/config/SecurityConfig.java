package com.coffeediary.config;

import com.coffeediary.service.AppUserDetailsService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.annotation.Order;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.ProviderManager;
import org.springframework.security.authentication.dao.DaoAuthenticationProvider;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.provisioning.InMemoryUserDetailsManager;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.context.HttpSessionSecurityContextRepository;
import org.springframework.security.web.context.SecurityContextRepository;

import org.springframework.session.jdbc.config.annotation.web.http.EnableJdbcHttpSession;
import org.springframework.session.web.http.CookieSerializer;
import org.springframework.session.web.http.DefaultCookieSerializer;

import java.io.IOException;
import java.util.Map;

@Configuration
@EnableWebSecurity
@EnableJdbcHttpSession(maxInactiveIntervalInSeconds = 7_776_000)
@RequiredArgsConstructor
public class SecurityConfig {

    private final AppUserDetailsService appUserDetailsService;

    @Value("${actuator.username}")
    private String actuatorUsername;

    @Value("${actuator.password}")
    private String actuatorPassword;

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public SecurityContextRepository securityContextRepository() {
        return new HttpSessionSecurityContextRepository();
    }

    @Bean
    public AuthenticationManager authenticationManager() {
        DaoAuthenticationProvider provider = new DaoAuthenticationProvider(appUserDetailsService);
        provider.setPasswordEncoder(passwordEncoder());
        return new ProviderManager(provider);
    }

    @Bean
    @Order(1)
    public SecurityFilterChain actuatorSecurityFilterChain(HttpSecurity http) throws Exception {
        UserDetailsService actuatorUserDetailsService = new InMemoryUserDetailsManager(
                User.builder()
                        .username(actuatorUsername)
                        .password(passwordEncoder().encode(actuatorPassword))
                        .roles("ACTUATOR")
                        .build()
        );

        http
                .securityMatcher("/actuator/**")
                .csrf(AbstractHttpConfigurer::disable)
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/actuator/health", "/actuator/info").permitAll()
                        .anyRequest().hasRole("ACTUATOR")
                )
                .httpBasic(basic -> basic
                        .realmName("Actuator")
                )
                .userDetailsService(actuatorUserDetailsService);

        return http.build();
    }

    @Bean
    @Order(2)
    public SecurityFilterChain apiSecurityFilterChain(HttpSecurity http) throws Exception {
        http
                .csrf(AbstractHttpConfigurer::disable)
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers("/api/auth/**").permitAll()
                        .requestMatchers("/api/**").authenticated()
                        .anyRequest().permitAll()
                )
                .securityContext(sc -> sc
                        .securityContextRepository(securityContextRepository())
                )
                .exceptionHandling(ex -> ex
                        .authenticationEntryPoint(jsonAuthenticationEntryPoint())
                )
                .headers(headers -> headers
                        .contentTypeOptions(contentTypeOptions -> {})
                        .frameOptions(frameOptions -> frameOptions.deny())
                        .httpStrictTransportSecurity(hsts -> hsts
                                .includeSubDomains(true)
                                .maxAgeInSeconds(31536000)
                        )
                        .contentSecurityPolicy(csp -> csp
                                .policyDirectives("default-src 'self'")
                        )
                )
                .userDetailsService(appUserDetailsService);

        return http.build();
    }

    @Bean
    public CookieSerializer cookieSerializer(
            @Value("${session.cookie.secure:false}") boolean secure) {
        DefaultCookieSerializer serializer = new DefaultCookieSerializer();
        serializer.setCookieMaxAge((int) java.time.Duration.ofDays(90).toSeconds());
        serializer.setUseSecureCookie(secure);
        serializer.setSameSite("Lax");
        return serializer;
    }

    @Bean
    public AuthenticationEntryPoint jsonAuthenticationEntryPoint() {
        return (HttpServletRequest request, HttpServletResponse response, AuthenticationException authException) -> {
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            response.setContentType("application/json");
            ObjectMapper mapper = new ObjectMapper();
            mapper.writeValue(response.getOutputStream(), Map.of(
                    "status", 401,
                    "error", "Unauthorized",
                    "message", "Authentication required"
            ));
        };
    }
}
