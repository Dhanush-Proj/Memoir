const API_URL = 'http://localhost:8080/blogs';

const blogForm = document.getElementById('blog-form');
const blogIdInput = document.getElementById('blog-id');
const blogTitleInput = document.getElementById('blog-title');
const blogDateInput = document.getElementById('blog-date');
const blogContentInput = document.getElementById('blog-content');
const blogImageInput = document.getElementById('blog-image');
const submitBtn = document.getElementById('submit-btn');
const blogsList = document.getElementById('blogs-list');

const addBlogBtn = document.getElementById('add-blog-btn');
const cancelBtn = document.getElementById('cancel-btn');

// --- New Dark Mode Elements ---
const modeToggleBtn = document.getElementById('mode-toggle');
const body = document.body;

// Function to fetch and display all blogs
async function fetchBlogs() {
    try {
        const response = await fetch(API_URL);
        const blogs = await response.json();
        
        // blogsList.innerHTML = '<h2>Recent Posts</h2>'; 
        
        if (blogs.length === 0) {
            blogsList.innerHTML += '<p>No blogs found. Add one above!</p>';
            return;
        }

        blogs.forEach(blog => {
            const blogPost = document.createElement('div');
            blogPost.className = 'blog-post';
            blogPost.innerHTML = `
                <h3>${blog.title}</h3>
                <p><strong>Date:</strong> ${blog.date}</p>
                <p>${blog.content}</p>
                ${blog.image ? `<img src="${blog.image}" alt="${blog.title}" style="max-width: 50%; height: 50%;  ">` : ''}
                <div class="actions">
                    <button class="edit-btn" onclick="editBlog('${blog.id}', '${blog.title}', '${blog.date}', '${blog.content}', '${blog.image}')">Edit</button>
                    <button class="delete-btn" onclick="deleteBlog('${blog.id}')">Delete</button>
                </div>
            `;
            blogsList.appendChild(blogPost);
        });
    } catch (error) {
        console.error('Error fetching blogs:', error);
    }
}

// Function to handle form submission (add/edit)
blogForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const id = blogIdInput.value;
    const blogData = {
        date: blogDateInput.value,
        title: blogTitleInput.value,
        content: blogContentInput.value,
        image: blogImageInput.value,
    };
    if (id) {
        await fetch(`${API_URL}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(blogData),
        });
    } else {
        await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(blogData),
        });
    }
    blogForm.style.display = 'none';
    blogsList.style.display = 'block';
    addBlogBtn.style.display = 'block';
    blogForm.reset();
    blogIdInput.value = '';
    submitBtn.textContent = 'Add Blog';
    fetchBlogs();
});

// Function to populate the form for editing
function editBlog(id, title, date, content, image) {
    addBlogBtn.style.display = 'none';
    blogForm.style.display = 'block';
    blogsList.style.display = 'none';
    blogIdInput.value = id;
    blogTitleInput.value = title;
    blogDateInput.value = date;
    blogContentInput.value = content;
    blogImageInput.value = image;
    submitBtn.textContent = 'Update Blog';
    window.scrollTo(0, 0);
}

// Function to delete a blog
async function deleteBlog(id) {
    if (confirm('Are you sure you want to delete this blog post?')) {
        await fetch(`${API_URL}/${id}`, {
            method: 'DELETE',
        });
        fetchBlogs();
    }
}

// Event listener for "Add New Blog" button
addBlogBtn.addEventListener('click', () => {
    addBlogBtn.style.display = 'none';
    blogForm.style.display = 'block';
    blogsList.style.display = 'none';
    blogForm.reset();
    blogIdInput.value = '';
    submitBtn.textContent = 'Add Blog';
});

// Event listener for "Cancel" button
cancelBtn.addEventListener('click', () => {
    blogForm.style.display = 'none';
    blogsList.style.display = 'block';
    addBlogBtn.style.display = 'block';
});

// --- New Dark Mode Logic ---
modeToggleBtn.addEventListener('click', () => {
    body.classList.toggle('dark-mode');
});

// Initial call to fetch blogs on page load
fetchBlogs();

