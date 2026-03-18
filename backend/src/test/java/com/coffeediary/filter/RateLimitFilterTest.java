package com.coffeediary.filter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;

import java.io.IOException;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.*;

class RateLimitFilterTest {

    private RateLimitFilter filter;
    private FilterChain filterChain;

    @BeforeEach
    void setUp() {
        filter = new RateLimitFilter(10.0);
        filterChain = mock(FilterChain.class);
    }

    @Test
    void allowsRequestUnderLimit() throws ServletException, IOException {
        MockHttpServletRequest request = new MockHttpServletRequest();
        request.setRemoteAddr("192.168.1.1");
        MockHttpServletResponse response = new MockHttpServletResponse();

        filter.doFilterInternal(request, response, filterChain);

        assertThat(response.getStatus()).isEqualTo(200);
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void usesXForwardedForHeader() throws ServletException, IOException {
        MockHttpServletRequest request = new MockHttpServletRequest();
        request.addHeader("X-Forwarded-For", "10.0.0.1, 192.168.1.1");
        request.setRemoteAddr("127.0.0.1");
        MockHttpServletResponse response = new MockHttpServletResponse();

        filter.doFilterInternal(request, response, filterChain);

        verify(filterChain).doFilter(request, response);
    }

    @Test
    void returns429WhenRateLimitExceeded() throws ServletException, IOException {
        // Exhaust the rate limiter for a specific IP
        String testIp = "10.99.99.99";

        MockHttpServletResponse lastResponse = null;
        boolean rateLimited = false;

        // Send requests until we get rate limited (limiter allows 10/second)
        for (int i = 0; i < 20; i++) {
            MockHttpServletRequest request = new MockHttpServletRequest();
            request.setRemoteAddr(testIp);
            lastResponse = new MockHttpServletResponse();
            filter.doFilterInternal(request, lastResponse, filterChain);

            if (lastResponse.getStatus() == 429) {
                rateLimited = true;
                break;
            }
        }

        assertThat(rateLimited).isTrue();
        assertThat(lastResponse.getStatus()).isEqualTo(429);
        assertThat(lastResponse.getContentType()).isEqualTo("application/json");
        assertThat(lastResponse.getContentAsString()).contains("Rate limit exceeded");
    }

    @Test
    void differentIpsHaveSeparateLimits() throws ServletException, IOException {
        MockHttpServletRequest request1 = new MockHttpServletRequest();
        request1.setRemoteAddr("10.0.0.1");
        MockHttpServletResponse response1 = new MockHttpServletResponse();

        MockHttpServletRequest request2 = new MockHttpServletRequest();
        request2.setRemoteAddr("10.0.0.2");
        MockHttpServletResponse response2 = new MockHttpServletResponse();

        filter.doFilterInternal(request1, response1, filterChain);
        filter.doFilterInternal(request2, response2, filterChain);

        assertThat(response1.getStatus()).isEqualTo(200);
        assertThat(response2.getStatus()).isEqualTo(200);
        verify(filterChain, times(2)).doFilter(any(), any());
    }
}
