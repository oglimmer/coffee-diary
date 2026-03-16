package com.coffeediary.repository;

import com.coffeediary.entity.DiaryEntry;
import org.springframework.data.jpa.domain.Specification;

import java.time.LocalDateTime;

public final class DiaryEntrySpecification {

    private DiaryEntrySpecification() {
    }

    public static Specification<DiaryEntry> belongsToUser(Long userId) {
        return (root, query, cb) -> cb.equal(root.get("user").get("id"), userId);
    }

    public static Specification<DiaryEntry> hasCoffeeId(Long coffeeId) {
        return (root, query, cb) -> cb.equal(root.get("coffee").get("id"), coffeeId);
    }

    public static Specification<DiaryEntry> hasSieveId(Long sieveId) {
        return (root, query, cb) -> cb.equal(root.get("sieve").get("id"), sieveId);
    }

    public static Specification<DiaryEntry> dateFrom(LocalDateTime from) {
        return (root, query, cb) -> cb.greaterThanOrEqualTo(root.get("dateTime"), from);
    }

    public static Specification<DiaryEntry> dateTo(LocalDateTime to) {
        return (root, query, cb) -> cb.lessThanOrEqualTo(root.get("dateTime"), to);
    }

    public static Specification<DiaryEntry> ratingMin(Integer ratingMin) {
        return (root, query, cb) -> cb.greaterThanOrEqualTo(root.get("rating"), ratingMin);
    }
}
