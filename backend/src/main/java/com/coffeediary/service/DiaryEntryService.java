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
import com.coffeediary.repository.DiaryEntrySpecification;
import com.coffeediary.repository.SieveRepository;
import com.coffeediary.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

@Slf4j
@Service
@RequiredArgsConstructor
public class DiaryEntryService {

    private final DiaryEntryRepository diaryEntryRepository;
    private final UserRepository userRepository;
    private final SieveRepository sieveRepository;
    private final CoffeeRepository coffeeRepository;

    @Transactional(readOnly = true)
    public PageResponse<DiaryEntryResponse> findAll(Long userId, Long coffeeId, Long sieveId,
                                                     LocalDateTime dateFrom, LocalDateTime dateTo,
                                                     Integer ratingMin, Pageable pageable) {
        Specification<DiaryEntry> spec = Specification.where(DiaryEntrySpecification.belongsToUser(userId));

        if (coffeeId != null) {
            spec = spec.and(DiaryEntrySpecification.hasCoffeeId(coffeeId));
        }
        if (sieveId != null) {
            spec = spec.and(DiaryEntrySpecification.hasSieveId(sieveId));
        }
        if (dateFrom != null) {
            spec = spec.and(DiaryEntrySpecification.dateFrom(dateFrom));
        }
        if (dateTo != null) {
            spec = spec.and(DiaryEntrySpecification.dateTo(dateTo));
        }
        if (ratingMin != null) {
            spec = spec.and(DiaryEntrySpecification.ratingMin(ratingMin));
        }

        Page<DiaryEntry> page = diaryEntryRepository.findAll(spec, pageable);

        return PageResponse.<DiaryEntryResponse>builder()
                .content(page.getContent().stream().map(this::toResponse).toList())
                .totalElements(page.getTotalElements())
                .totalPages(page.getTotalPages())
                .number(page.getNumber())
                .size(page.getSize())
                .build();
    }

    @Transactional(readOnly = true)
    public DiaryEntryResponse findById(Long userId, Long entryId) {
        DiaryEntry entry = diaryEntryRepository.findById(entryId)
                .orElseThrow(() -> new IllegalArgumentException("Diary entry not found"));

        if (!entry.getUser().getId().equals(userId)) {
            throw new SecurityException("Not authorized to view this entry");
        }

        return toResponse(entry);
    }

    @Transactional
    public DiaryEntryResponse create(Long userId, DiaryEntryRequest request) {
        User user = userRepository.getReferenceById(userId);

        DiaryEntry entry = DiaryEntry.builder()
                .user(user)
                .dateTime(request.getDateTime())
                .temperature(request.getTemperature() != null ? request.getTemperature() : 93)
                .grindSize(request.getGrindSize())
                .inputWeight(request.getInputWeight())
                .outputWeight(request.getOutputWeight())
                .timeSeconds(request.getTimeSeconds())
                .rating(request.getRating())
                .notes(request.getNotes())
                .build();

        if (request.getSieveId() != null) {
            Sieve sieve = sieveRepository.findById(request.getSieveId())
                    .orElseThrow(() -> new IllegalArgumentException("Sieve not found"));
            if (!sieve.getUser().getId().equals(userId)) {
                throw new IllegalArgumentException("Sieve does not belong to user");
            }
            entry.setSieve(sieve);
        }

        if (request.getCoffeeId() != null) {
            Coffee coffee = coffeeRepository.findById(request.getCoffeeId())
                    .orElseThrow(() -> new IllegalArgumentException("Coffee not found"));
            if (!coffee.getUser().getId().equals(userId)) {
                throw new IllegalArgumentException("Coffee does not belong to user");
            }
            entry.setCoffee(coffee);
        }

        entry = diaryEntryRepository.save(entry);
        log.info("Created diary entry {} for user {}", entry.getId(), userId);

        return toResponse(entry);
    }

    @Transactional
    public DiaryEntryResponse update(Long userId, Long entryId, DiaryEntryRequest request) {
        DiaryEntry entry = diaryEntryRepository.findById(entryId)
                .orElseThrow(() -> new IllegalArgumentException("Diary entry not found"));

        if (!entry.getUser().getId().equals(userId)) {
            throw new SecurityException("Not authorized to update this entry");
        }

        entry.setDateTime(request.getDateTime());
        entry.setTemperature(request.getTemperature() != null ? request.getTemperature() : 93);
        entry.setGrindSize(request.getGrindSize());
        entry.setInputWeight(request.getInputWeight());
        entry.setOutputWeight(request.getOutputWeight());
        entry.setTimeSeconds(request.getTimeSeconds());
        entry.setRating(request.getRating());
        entry.setNotes(request.getNotes());

        if (request.getSieveId() != null) {
            Sieve sieve = sieveRepository.findById(request.getSieveId())
                    .orElseThrow(() -> new IllegalArgumentException("Sieve not found"));
            if (!sieve.getUser().getId().equals(userId)) {
                throw new IllegalArgumentException("Sieve does not belong to user");
            }
            entry.setSieve(sieve);
        } else {
            entry.setSieve(null);
        }

        if (request.getCoffeeId() != null) {
            Coffee coffee = coffeeRepository.findById(request.getCoffeeId())
                    .orElseThrow(() -> new IllegalArgumentException("Coffee not found"));
            if (!coffee.getUser().getId().equals(userId)) {
                throw new IllegalArgumentException("Coffee does not belong to user");
            }
            entry.setCoffee(coffee);
        } else {
            entry.setCoffee(null);
        }

        entry = diaryEntryRepository.save(entry);
        log.info("Updated diary entry {} for user {}", entry.getId(), userId);

        return toResponse(entry);
    }

    @Transactional
    public void delete(Long userId, Long entryId) {
        DiaryEntry entry = diaryEntryRepository.findById(entryId)
                .orElseThrow(() -> new IllegalArgumentException("Diary entry not found"));

        if (!entry.getUser().getId().equals(userId)) {
            throw new SecurityException("Not authorized to delete this entry");
        }

        diaryEntryRepository.delete(entry);
        log.info("Deleted diary entry {} for user {}", entryId, userId);
    }

    private DiaryEntryResponse toResponse(DiaryEntry entry) {
        return DiaryEntryResponse.builder()
                .id(entry.getId())
                .userId(entry.getUser().getId())
                .dateTime(entry.getDateTime())
                .sieveId(entry.getSieve() != null ? entry.getSieve().getId() : null)
                .sieveName(entry.getSieve() != null ? entry.getSieve().getName() : null)
                .temperature(entry.getTemperature())
                .coffeeId(entry.getCoffee() != null ? entry.getCoffee().getId() : null)
                .coffeeName(entry.getCoffee() != null ? entry.getCoffee().getName() : null)
                .grindSize(entry.getGrindSize())
                .inputWeight(entry.getInputWeight())
                .outputWeight(entry.getOutputWeight())
                .timeSeconds(entry.getTimeSeconds())
                .rating(entry.getRating())
                .notes(entry.getNotes())
                .build();
    }
}
