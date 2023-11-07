might be useful: 
https://codevoweb.com/golang-and-gorm-user-registration-email-verification/
crud:
https://www.youtube.com/watch?v=lf_kiH_NPvM

jwt & auth:
https://www.youtube.com/watch?v=ma7rUS_vW9M

architecture: 
https://www.youtube.com/watch?v=EqniGcAijDI
## Igor's suggestions
### Simple REST API for a To-Do List

1. Basic Features:

    a. CRUD Operations:
    X Create tasks. 
    X Read tasks (single & all).
    X Update tasks.
    X Delete tasks.
    X uuid implemented.

2. Intermediate Features:

    a. User Authentication & Authorization:
    - Sign up & login functionalities. <!-- sorta kinda  -->
    - JWT or cookie-based sessions.
    - Role-based access (e.g., admin, user).

    b. Persistence:
    - Integrate a relational database like PostgreSQL or SQLite.
    X Use an ORM (Object-Relational Mapping) tool like GORM.

    c. Middleware:
    - Logging incoming requests.
    - Rate limiting.
    - CORS (Cross-Origin Resource Sharing) settings.

    d. Validation & Error Handling:
    - Input validation for task data.
    - Consistent error response structure.

3. Advanced Features:

    a. Task Enhancements:
    - Categorize tasks.
    - Prioritize tasks.
    - Set due dates and reminders.

    b. Real-time Features:
    - Implement WebSocket to update tasks in real-time.

    c. Pagination & Filtering:
    - Allow large sets of tasks to be viewed in pages.
    - Filter tasks by different criteria (e.g., completed, due soon).

    d. API Versioning:
    - Maintain multiple versions of the API for backward compatibility.

    e. File Attachments:
    - Attach files or images to tasks.
    - Store files in cloud storage (like AWS S3) or locally.

4. Expert Features:

    a. Third-party Integration:
    - Sync tasks with external tools (e.g., Google Calendar).
    - Allow sharing tasks on social media.

    b. Analytics:
    - Track how often tasks are added, completed, or updated.
    - Generate reports on user activity.

    c. Testing & Deployment:
    - Write unit and integration tests using Go's testing framework.
    - Set up Continuous Integration and Continuous Deployment (CI/CD).
    - Deploy the API to cloud providers (e.g., AWS, Google Cloud).

    d. Documentation & SDKs:
    - Use tools like Swagger for API documentation.
    - Provide SDKs or client libraries in other languages for easier API consumption.

5. Monetization & Scaling (Optional):
    
    a. Premium Features:
    - Offer advanced features (e.g., collaboration, advanced analytics) for a subscription fee.

    b. Scaling:
    - Implement caching using tools like Redis.
    - Optimize database queries.
    - Load balancing.
