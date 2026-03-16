package com.coffeediary.controller;

import com.coffeediary.dto.CoffeeRequest;
import com.coffeediary.dto.CoffeeResponse;
import com.coffeediary.service.AppUserDetails;
import com.coffeediary.service.CoffeeService;
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
@RequestMapping("/api/coffees")
@RequiredArgsConstructor
public class CoffeeController {

    private final CoffeeService coffeeService;

    @GetMapping
    public ResponseEntity<List<CoffeeResponse>> findAll(@AuthenticationPrincipal AppUserDetails userDetails) {
        return ResponseEntity.ok(coffeeService.findAllByUser(userDetails.getId()));
    }

    @PostMapping
    public ResponseEntity<CoffeeResponse> create(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @Valid @RequestBody CoffeeRequest request) {

        CoffeeResponse response = coffeeService.create(userDetails.getId(), request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(
            @AuthenticationPrincipal AppUserDetails userDetails,
            @PathVariable Long id) {

        coffeeService.delete(userDetails.getId(), id);
        return ResponseEntity.noContent().build();
    }
}
