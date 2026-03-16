package com.coffeediary.repository;

import com.coffeediary.entity.Sieve;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface SieveRepository extends JpaRepository<Sieve, Long> {

    List<Sieve> findAllByUserId(Long userId);
}
