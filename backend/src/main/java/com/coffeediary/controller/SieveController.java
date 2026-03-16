package com.coffeediary.controller;

import com.coffeediary.dto.SieveRequest;
import com.coffeediary.dto.SieveResponse;
import com.coffeediary.service.AppUserDetails;
import com.coffeediary.service.SieveService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@RestController
@RequestMapping("/api/sieves")
@RequiredArgsConstructor
public class SieveController {

    private final SieveService sieveService;

    @GetMapping
    public ResponseEntity<List<SieveResponse>> findAll(@AuthenticationPrincipal AppUserDetails userDetails) {
        return ResponseEntity.ok(sieveService.findAllByUser(userDetails.getId()));
    }

    @PostMapping
    public ResponseEntity<SieveResponse> create(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @Valid @RequestBody SieveRequest request) {

        SieveResponse response = sieveService.create(userDetails.getId(), request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @PathVariable Long id) {

        sieveService.delete(userDetails.getId(), id);
        return ResponseEntity.noContent().build();
    }
}
