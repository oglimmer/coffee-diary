package com.coffeediary.service;

import com.coffeediary.dto.SieveRequest;
import com.coffeediary.dto.SieveResponse;
import com.coffeediary.entity.Sieve;
import com.coffeediary.entity.User;
import com.coffeediary.repository.SieveRepository;
import com.coffeediary.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class SieveService {

    private final SieveRepository sieveRepository;
    private final UserRepository userRepository;

    @Transactional(readOnly = true)
    public List<SieveResponse> findAllByUser(Long userId) {
        return sieveRepository.findAllByUserId(userId).stream()
                .map(this::toResponse)
                .toList();
    }

    @Transactional
    public SieveResponse create(Long userId, SieveRequest request) {
        User user = userRepository.getReferenceById(userId);

        Sieve sieve = Sieve.builder()
                .name(request.getName())
                .user(user)
                .build();

        sieve = sieveRepository.save(sieve);
        log.info("Created sieve {} for user {}", sieve.getId(), userId);

        return toResponse(sieve);
    }

    @Transactional
    public void delete(Long userId, Long sieveId) {
        Sieve sieve = sieveRepository.findById(sieveId)
                .orElseThrow(() -> new IllegalArgumentException("Sieve not found"));

        if (!sieve.getUser().getId().equals(userId)) {
            throw new SecurityException("Not authorized to delete this sieve");
        }

        sieveRepository.delete(sieve);
        log.info("Deleted sieve {} for user {}", sieveId, userId);
    }

    private SieveResponse toResponse(Sieve sieve) {
        return SieveResponse.builder()
                .id(sieve.getId())
                .name(sieve.getName())
                .build();
    }
}
