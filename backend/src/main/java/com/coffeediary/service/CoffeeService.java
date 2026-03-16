package com.coffeediary.service;

import com.coffeediary.dto.CoffeeRequest;
import com.coffeediary.dto.CoffeeResponse;
import com.coffeediary.entity.Coffee;
import com.coffeediary.entity.User;
import com.coffeediary.repository.CoffeeRepository;
import com.coffeediary.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class CoffeeService {

    private final CoffeeRepository coffeeRepository;
    private final UserRepository userRepository;

    @Transactional(readOnly = true)
    public List<CoffeeResponse> findAllByUser(Long userId) {
        return coffeeRepository.findAllByUserId(userId).stream()
                .map(this::toResponse)
                .toList();
    }

    @Transactional
    public CoffeeResponse create(Long userId, CoffeeRequest request) {
        User user = userRepository.getReferenceById(userId);

        Coffee coffee = Coffee.builder()
                .name(request.getName())
                .user(user)
                .build();

        coffee = coffeeRepository.save(coffee);
        log.info("Created coffee {} for user {}", coffee.getId(), userId);

        return toResponse(coffee);
    }

    @Transactional
    public void delete(Long userId, Long coffeeId) {
        Coffee coffee = coffeeRepository.findById(coffeeId)
                .orElseThrow(() -> new IllegalArgumentException("Coffee not found"));

        if (!coffee.getUser().getId().equals(userId)) {
            throw new SecurityException("Not authorized to delete this coffee");
        }

        coffeeRepository.delete(coffee);
        log.info("Deleted coffee {} for user {}", coffeeId, userId);
    }

    private CoffeeResponse toResponse(Coffee coffee) {
        return CoffeeResponse.builder()
                .id(coffee.getId())
                .name(coffee.getName())
                .build();
    }
}
