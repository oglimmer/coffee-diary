package com.coffeediary.controller;

import com.coffeediary.dto.DiaryEntryRequest;
import com.coffeediary.dto.DiaryEntryResponse;
import com.coffeediary.dto.PageResponse;
import com.coffeediary.service.AppUserDetails;
import com.coffeediary.service.DiaryEntryService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;

@Slf4j
@RestController
@RequestMapping("/api/diary-entries")
@RequiredArgsConstructor
public class DiaryEntryController {

    private final DiaryEntryService diaryEntryService;

    @GetMapping
    public ResponseEntity<PageResponse<DiaryEntryResponse>> findAll(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @RequestParam(required = false) Long coffeeId,
            @RequestParam(required = false) Long sieveId,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime dateFrom,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime dateTo,
            @RequestParam(required = false) Integer ratingMin,
            @PageableDefault(size = 20, sort = "dateTime") Pageable pageable) {

        PageResponse<DiaryEntryResponse> response = diaryEntryService.findAll(
                userDetails.getId(), coffeeId, sieveId, dateFrom, dateTo, ratingMin, pageable);

        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<DiaryEntryResponse> findById(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @PathVariable Long id) {

        DiaryEntryResponse response = diaryEntryService.findById(userDetails.getId(), id);
        return ResponseEntity.ok(response);
    }

    @PostMapping
    public ResponseEntity<DiaryEntryResponse> create(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @Valid @RequestBody DiaryEntryRequest request) {

        DiaryEntryResponse response = diaryEntryService.create(userDetails.getId(), request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @PutMapping("/{id}")
    public ResponseEntity<DiaryEntryResponse> update(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @PathVariable Long id,
            @Valid @RequestBody DiaryEntryRequest request) {

        DiaryEntryResponse response = diaryEntryService.update(userDetails.getId(), id, request);
        return ResponseEntity.ok(response);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @PathVariable Long id) {

        diaryEntryService.delete(userDetails.getId(), id);
        return ResponseEntity.noContent().build();
    }
}
