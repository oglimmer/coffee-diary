package com.coffeediary.filter;

import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.LoadingCache;
import com.google.common.util.concurrent.RateLimiter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

@Component
@SuppressWarnings("UnstableApiUsage")
public class RateLimitFilter extends OncePerRequestFilter {

    private final LoadingCache<String, RateLimiter> rateLimiters;

    public RateLimitFilter(@Value("${rate-limit.requests-per-second:10.0}") double requestsPerSecond) {
        this.rateLimiters = CacheBuilder.newBuilder()
                .expireAfterAccess(10, TimeUnit.MINUTES)
                .build(new CacheLoader<>() {
                    @Override
                    public RateLimiter load(String key) {
                        return RateLimiter.create(requestsPerSecond);
                    }
                });
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        String clientIp = getClientIp(request);
        RateLimiter limiter = rateLimiters.getUnchecked(clientIp);

        if (!limiter.tryAcquire()) {
            response.setStatus(HttpStatus.TOO_MANY_REQUESTS.value());
            response.setContentType("application/json");
            response.getWriter().write("{\"status\":429,\"error\":\"Too Many Requests\",\"message\":\"Rate limit exceeded\"}");
            return;
        }

        filterChain.doFilter(request, response);
    }

    private String getClientIp(HttpServletRequest request) {
        String xForwardedFor = request.getHeader("X-Forwarded-For");
        if (xForwardedFor != null && !xForwardedFor.isEmpty()) {
            return xForwardedFor.split(",")[0].trim();
        }
        return request.getRemoteAddr();
    }
}
