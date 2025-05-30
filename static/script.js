// API Configuration
const API_BASE_URL = 'http://0.0.0.0:8080/api/v1';

// DOM Elements
const userForm = document.getElementById('user-form');
const usersContainer = document.getElementById('users-container');
const refreshBtn = document.getElementById('refresh-btn');
const editModal = document.getElementById('edit-modal');
const editForm = document.getElementById('edit-form');
const statusIndicator = document.getElementById('api-status');
const statusText = document.getElementById('status-text');
const messageContainer = document.getElementById('message-container');

// Initialize the app
document.addEventListener('DOMContentLoaded', function() {
    checkApiHealth();
    loadUsers();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    userForm.addEventListener('submit', handleCreateUser);
    editForm.addEventListener('submit', handleUpdateUser);
    refreshBtn.addEventListener('click', loadUsers);
    
    // Modal close events
    document.querySelector('.close').addEventListener('click', closeModal);
    window.addEventListener('click', function(event) {
        if (event.target === editModal) {
            closeModal();
        }
    });
}

// Check API health
async function checkApiHealth() {
    try {
        statusText.textContent = 'Checking API...';
        statusIndicator.className = 'status-indicator checking';
        
        const response = await fetch(`${API_BASE_URL}/health`);
        if (response.ok) {
            statusText.textContent = 'API is healthy';
            statusIndicator.className = 'status-indicator healthy';
        } else {
            throw new Error('API returned error status');
        }
    } catch (error) {
        statusText.textContent = 'API is not available';
        statusIndicator.className = 'status-indicator unhealthy';
        showMessage('Unable to connect to API. Make sure the server is running on port 8080.', 'error');
    }
}

// Load all users
async function loadUsers() {
    try {
        usersContainer.innerHTML = '<div class="loading">Loading users...</div>';
        
        const response = await fetch(`${API_BASE_URL}/users`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const users = await response.json();
        displayUsers(users);
    } catch (error) {
        console.error('Error loading users:', error);
        usersContainer.innerHTML = '<div class="loading">Failed to load users. Please check if the server is running.</div>';
        showMessage('Failed to load users: ' + error.message, 'error');
    }
}

// Display users in the UI
function displayUsers(users) {    
    if (!users || users.length === 0) {
        usersContainer.innerHTML = `
            <div class="empty-state">
                <h3>No users found</h3>
                <p>Add your first user using the form above.</p>
            </div>
        `;
        return;
    }

    const usersHTML = users.map(user => `
        <div class="user-card">
            <div class="user-header">
                <div class="user-id">ID: ${user.id}</div>
            </div>
            <div class="user-info">
                <h3>${escapeHtml(user.name)}</h3>
                <p><strong>Email:</strong> ${escapeHtml(user.email)}</p>
                <p><strong>Created:</strong> ${formatDate(user.created)}</p>
            </div>
            <div class="user-actions">
                <button class="btn btn-edit" onclick="editUser(${user.id})">Edit</button>
                <button class="btn btn-danger" onclick="deleteUser(${user.id})">Delete</button>
            </div>
        </div>
    `).join('');

    usersContainer.innerHTML = `<div class="users-grid">${usersHTML}</div>`;
}

// Handle create user form submission
async function handleCreateUser(event) {
    event.preventDefault();
    
    const formData = new FormData(userForm);
    const userData = {
        name: formData.get('name').trim(),
        email: formData.get('email').trim()
    };

    if (!userData.name || !userData.email) {
        showMessage('Please fill in all fields', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create user');
        }

        const newUser = await response.json();
        showMessage(`User "${newUser.name}" created successfully!`, 'success');
        userForm.reset();
        loadUsers();
    } catch (error) {
        showMessage(`Failed to create user: ${error.message}`, 'error');
        console.error('Error creating user:', error);
    }
}

// Edit user
async function editUser(userId) {
    try {
        const response = await fetch(`${API_BASE_URL}/users/${userId}`);
        if (!response.ok) {
            throw new Error('Failed to load user data');
        }

        const user = await response.json();
        
        // Populate edit form
        document.getElementById('edit-id').value = user.id;
        document.getElementById('edit-name').value = user.name;
        document.getElementById('edit-email').value = user.email;
        
        // Show modal
        editModal.style.display = 'block';
    } catch (error) {
        showMessage(`Failed to load user data: ${error.message}`, 'error');
        console.error('Error loading user for edit:', error);
    }
}

// Handle update user form submission
async function handleUpdateUser(event) {
    event.preventDefault();
    
    const userId = document.getElementById('edit-id').value;
    const formData = new FormData(editForm);
    const userData = {
        name: formData.get('name').trim(),
        email: formData.get('email').trim()
    };

    if (!userData.name || !userData.email) {
        showMessage('Please fill in all fields', 'error');
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to update user');
        }

        const updatedUser = await response.json();
        showMessage(`User "${updatedUser.name}" updated successfully!`, 'success');
        closeModal();
        loadUsers();
    } catch (error) {
        showMessage(`Failed to update user: ${error.message}`, 'error');
        console.error('Error updating user:', error);
    }
}

// Delete user
async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Failed to delete user');
        }

        showMessage('User deleted successfully!', 'success');
        loadUsers();
    } catch (error) {
        showMessage(`Failed to delete user: ${error.message}`, 'error');
        console.error('Error deleting user:', error);
    }
}

// Close modal
function closeModal() {
    editModal.style.display = 'none';
    editForm.reset();
}

// Show message to user
function showMessage(text, type = 'success') {
    const message = document.createElement('div');
    message.className = `message ${type}`;
    message.textContent = text;
    
    messageContainer.appendChild(message);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (message.parentNode) {
            message.parentNode.removeChild(message);
        }
    }, 5000);
}

// Utility functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    } catch (error) {
        return dateString;
    }
}
