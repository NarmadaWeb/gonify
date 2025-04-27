// Main JavaScript file for Minification Test Project

// DOM Elements
const ctaButton = document.getElementById('ctaButton');
const dataContainer = document.getElementById('dataContainer');

// Configuration Object
const appConfig = {
    apiEndpoint: 'https://api.example.com/data',
    maxRetries: 3,
    timeout: 5000,
    features: {
        darkMode: true,
        animations: true,
        analytics: false
    }
};

// Sample Data
const sampleData = [
    { id: 1, name: 'HTML Minifier', version: '1.2.0', downloads: 12500 },
    { id: 2, name: 'CSS Optimizer', version: '2.1.3', downloads: 8900 },
    { id: 3, name: 'JS Compressor', version: '3.0.5', downloads: 21500 },
    { id: 4, name: 'Image Compressor', version: '1.0.8', downloads: 17800 }
];

// Utility Functions
function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function createElement(tag, attributes = {}, textContent = '') {
    const element = document.createElement(tag);
    Object.keys(attributes).forEach(key => {
        element.setAttribute(key, attributes[key]);
    });
    if (textContent) element.textContent = textContent;
    return element;
}

// Event Handlers
function handleCtaClick(event) {
    console.log('CTA button clicked!');
    fetchData();
    event.target.disabled = true;
    event.target.textContent = 'Loading...';
}

async function fetchData() {
    try {
        // Simulate API call
        console.log('Fetching data from:', appConfig.apiEndpoint);

        // Using timeout to simulate network request
        await new Promise(resolve => setTimeout(resolve, 1500));

        // Process sample data instead of real API call
        renderData(sampleData);

        // Re-enable button
        ctaButton.disabled = false;
        ctaButton.textContent = 'Refresh Data';
    } catch (error) {
        console.error('Error fetching data:', error);
        ctaButton.disabled = false;
        ctaButton.textContent = 'Try Again';
    }
}

function renderData(data) {
    // Clear previous content
    dataContainer.innerHTML = '';

    // Create table element
    const table = createElement('table', { class: 'data-table' });
    const thead = createElement('thead');
    const tbody = createElement('tbody');

    // Create table header
    const headerRow = createElement('tr');
    ['ID', 'Tool Name', 'Version', 'Downloads'].forEach(text => {
        headerRow.appendChild(createElement('th', {}, text));
    });
    thead.appendChild(headerRow);

    // Create table rows with data
    data.forEach(item => {
        const row = createElement('tr');
        row.appendChild(createElement('td', {}, item.id));
        row.appendChild(createElement('td', {}, item.name));
        row.appendChild(createElement('td', {}, item.version));
        row.appendChild(createElement('td', {}, formatNumber(item.downloads)));
        tbody.appendChild(row);
    });

    // Assemble table
    table.appendChild(thead);
    table.appendChild(tbody);
    dataContainer.appendChild(table);

    // Add some basic styles dynamically
    const style = createElement('style');
    style.textContent = `
        .data-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 2rem;
        }
        .data-table th, .data-table td {
            padding: 0.75rem;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .data-table th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        .data-table tr:hover {
            background-color: #f5f5f5;
        }
    `;
    document.head.appendChild(style);
}

// Initialize the application
function init() {
    console.log('Application initialized');
    ctaButton.addEventListener('click', handleCtaClick);

    // Check for features
    if (appConfig.features.darkMode) {
        console.log('Dark mode feature is enabled');
    }

    // Load some initial data
    setTimeout(() => {
        console.log('Loading initial data...');
        renderData(sampleData.slice(0, 2));
    }, 500);
}

// Start the app when DOM is loaded
document.addEventListener('DOMContentLoaded', init);
