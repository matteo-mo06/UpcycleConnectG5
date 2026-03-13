<?php

use Illuminate\Support\Facades\Route;
use Inertia\Inertia;

Route::get('/', function () {
    return Inertia::render('Home');
});

Route::prefix('admin')->name('admin.')->group(function () {
    Route::get('/dashboard', function () {
        return Inertia::render('Admin/Dashboard');
    })->name('dashboard');

    Route::get('/users', function () {
        return Inertia::render('Admin/Users');
    })->name('users');

    Route::get('/listings', function () {
        return Inertia::render('Admin/Listings');
    })->name('listings');

    Route::get('/events', function () {
        return Inertia::render('Admin/Events');
    })->name('events');

    Route::get('/reports', function () {
        return Inertia::render('Admin/Reports');
    })->name('reports');
});