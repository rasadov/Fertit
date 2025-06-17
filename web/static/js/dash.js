// Utility functions
function showAlert(message, type = 'success') {
    const alertsContainer = document.getElementById('alerts');
    const alert = document.createElement('div');
    alert.className = `alert alert-${type}`;
    alert.textContent = message;
    alertsContainer.appendChild(alert);

    setTimeout(() => {
        alert.remove();
    }, 5000);
}

function setLoading(form, loading = true) {
    if (loading) {
        form.classList.add('loading');
        const button = form.querySelector('button[type="submit"]');
        button.textContent = 'Processing...';
        button.disabled = true;
    } else {
        form.classList.remove('loading');
        const button = form.querySelector('button[type="submit"]');
        button.textContent = form.id === 'emailForm' ? 'Send Newsletter' : 'Add Subscriber';
        button.disabled = false;
    }
}

// Email form submission
document.getElementById('emailForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const form = this;
    setLoading(form, true);

    const formData = new FormData(form);

    try {
        const response = await fetch('/admin/send-email', {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            showAlert('Newsletter sent successfully!', 'success');
            form.reset();
        } else {
            const error = await response.json();
            showAlert(error.error || 'Failed to send newsletter', 'error');
        }
    } catch (error) {
        showAlert('Network error occurred', 'error');
    } finally {
        setLoading(form, false);
    }
});

// User form submission (using subscribe route)
document.getElementById('userForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const form = this;
    setLoading(form, true);

    const formData = new FormData(form);

    // Handle multiple categories - join them with commas for the subscribe route
    const selectedCategories = Array.from(document.getElementById('userCategory').selectedOptions)
        .map(option => option.value);

    if (selectedCategories.length > 0) {
        formData.set('categories', selectedCategories.join(','));
    } else {
        formData.set('categories', 'general');
    }

    try {
        const response = await fetch('/subscribe', {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            showAlert('Subscriber added successfully!', 'success');
            form.reset();
            // Refresh the page after 2 seconds to show the new subscriber
            setTimeout(() => {
                location.reload();
            }, 2000);
        } else {
            const error = await response.json();
            showAlert(error.error || 'Failed to add subscriber', 'error');
        }
    } catch (error) {
        showAlert('Network error occurred', 'error');
    } finally {
        setLoading(form, false);
    }
});

// Elements per page change
document.getElementById('elementsPerPage').addEventListener('change', function() {
    const elements = this.value;
    const url = new URL(window.location);
    url.searchParams.set('elements', elements);
    url.searchParams.set('page', '1'); // Reset to first page
    window.location.href = url.toString();
});

// Auto-resize textarea
document.getElementById('body').addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = this.scrollHeight + 'px';
});