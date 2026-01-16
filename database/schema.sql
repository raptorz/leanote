-- Leanote PostgreSQL Database Schema
-- This schema replaces MongoDB collections with PostgreSQL tables

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- User Tables
-- ============================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN DEFAULT FALSE,
    username VARCHAR(255) NOT NULL UNIQUE,
    username_raw VARCHAR(255),
    pwd VARCHAR(255),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    logo VARCHAR(500),
    theme VARCHAR(100),
    notebook_width INT DEFAULT 200,
    note_list_width INT DEFAULT 200,
    md_editor_width INT DEFAULT 0,
    left_is_min BOOLEAN DEFAULT FALSE,
    third_user_id VARCHAR(255),
    third_username VARCHAR(255),
    third_type INT DEFAULT 0,
    image_num INT DEFAULT 0,
    image_size BIGINT DEFAULT 0,
    attach_num INT DEFAULT 0,
    attach_size BIGINT DEFAULT 0,
    from_user_id UUID,
    account_type VARCHAR(50) DEFAULT 'normal',
    account_start_time TIMESTAMP,
    account_end_time TIMESTAMP,
    max_image_num INT,
    max_image_size BIGINT,
    max_attach_num INT,
    max_attach_size BIGINT,
    max_per_attach_size BIGINT,
    usn INT DEFAULT 0,
    full_sync_before TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_from_user_id ON users(from_user_id);

-- ============================================
-- Notebook Tables
-- ============================================

CREATE TABLE notebooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_notebook_id UUID REFERENCES notebooks(id) ON DELETE CASCADE,
    seq INT DEFAULT 0,
    title VARCHAR(255) NOT NULL,
    url_title VARCHAR(255),
    number_notes INT DEFAULT 0,
    is_trash BOOLEAN DEFAULT FALSE,
    is_blog BOOLEAN DEFAULT FALSE,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usn INT DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_notebooks_user_id ON notebooks(user_id);
CREATE INDEX idx_notebooks_parent_id ON notebooks(parent_notebook_id);
CREATE INDEX idx_notebooks_is_trash ON notebooks(is_trash);
CREATE INDEX idx_notebooks_is_deleted ON notebooks(is_deleted);

-- ============================================
-- Note Tables
-- ============================================

CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    notebook_id UUID NOT NULL REFERENCES notebooks(id) ON DELETE CASCADE,
    title VARCHAR(500),
    description TEXT,
    src VARCHAR(500),
    img_src VARCHAR(500),
    tags TEXT[],
    is_trash BOOLEAN DEFAULT FALSE,
    is_blog BOOLEAN DEFAULT FALSE,
    url_title VARCHAR(500),
    is_recommend BOOLEAN DEFAULT FALSE,
    is_top BOOLEAN DEFAULT FALSE,
    has_self_defined BOOLEAN DEFAULT FALSE,
    read_num INT DEFAULT 0,
    like_num INT DEFAULT 0,
    comment_num INT DEFAULT 0,
    is_markdown BOOLEAN DEFAULT FALSE,
    attach_num INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    recommend_time TIMESTAMP,
    public_time TIMESTAMP,
    updated_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    usn INT DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_notebook_id ON notes(notebook_id);
CREATE INDEX idx_notes_created_user_id ON notes(created_user_id);
CREATE INDEX idx_notes_is_trash ON notes(is_trash);
CREATE INDEX idx_notes_is_blog ON notes(is_blog);
CREATE INDEX idx_notes_is_deleted ON notes(is_deleted);
CREATE INDEX idx_notes_tags ON notes USING GIN(tags);

-- Note Content Table
CREATE TABLE note_contents (
    note_id UUID PRIMARY KEY REFERENCES notes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_blog BOOLEAN DEFAULT FALSE,
    content TEXT,
    abstract TEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_user_id UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_note_contents_user_id ON note_contents(user_id);
CREATE INDEX idx_note_contents_is_blog ON note_contents(is_blog);

-- Note Content History Table
CREATE TABLE note_content_histories (
    note_id UUID PRIMARY KEY REFERENCES notes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    histories JSONB
);

CREATE INDEX idx_note_content_histories_user_id ON note_content_histories(user_id);

-- ============================================
-- Tag Tables
-- ============================================

CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag VARCHAR(255) NOT NULL,
    usn INT DEFAULT 0,
    count INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, tag)
);

CREATE INDEX idx_tags_user_id ON tags(user_id);
CREATE INDEX idx_tags_is_deleted ON tags(is_deleted);

CREATE TABLE tag_counts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag VARCHAR(255) NOT NULL,
    is_blog BOOLEAN DEFAULT FALSE,
    count INT DEFAULT 0,
    UNIQUE(user_id, tag, is_blog)
);

CREATE INDEX idx_tag_counts_user_id ON tag_counts(user_id);
CREATE INDEX idx_tag_counts_is_blog ON tag_counts(is_blog);

-- ============================================
-- Attachment Tables
-- ============================================

CREATE TABLE attachs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    note_id UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    upload_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    size BIGINT NOT NULL,
    type VARCHAR(50),
    path VARCHAR(500),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attachs_note_id ON attachs(note_id);
CREATE INDEX idx_attachs_upload_user_id ON attachs(upload_user_id);

-- ============================================
-- Share Tables
-- ============================================

CREATE TABLE share_notebooks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    notebook_id UUID NOT NULL REFERENCES notebooks(id) ON DELETE CASCADE,
    permissions INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_share_notebooks_user_id ON share_notebooks(user_id);
CREATE INDEX idx_share_notebooks_to_user_id ON share_notebooks(to_user_id);
CREATE INDEX idx_share_notebooks_notebook_id ON share_notebooks(notebook_id);

CREATE TABLE share_notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    note_id UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    permissions INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_share_notes_user_id ON share_notes(user_id);
CREATE INDEX idx_share_notes_note_id ON share_notes(note_id);

-- ============================================
-- Group Tables
-- ============================================

CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    user_count INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_groups_user_id ON groups(user_id);
CREATE INDEX idx_groups_is_deleted ON groups(is_deleted);

CREATE TABLE group_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_group_users_group_id ON group_users(group_id);
CREATE INDEX idx_group_users_user_id ON group_users(user_id);

-- ============================================
-- Blog Tables
-- ============================================

CREATE TABLE blogs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500),
    subtitle VARCHAR(500),
    about_me TEXT,
    logo VARCHAR(500),
    domain VARCHAR(255),
    theme VARCHAR(100),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_blogs_user_id ON blogs(user_id);
CREATE INDEX idx_blogs_is_deleted ON blogs(is_deleted);

CREATE TABLE blog_singles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    note_id UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    title VARCHAR(500),
    slug VARCHAR(500),
    is_published BOOLEAN DEFAULT FALSE,
    published_time TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_blog_singles_note_id ON blog_singles(note_id);
CREATE INDEX idx_blog_singles_is_published ON blog_singles(is_published);

CREATE TABLE blog_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(blog_id, user_id)
);

CREATE INDEX idx_blog_likes_blog_id ON blog_likes(blog_id);
CREATE INDEX idx_blog_likes_user_id ON blog_likes(user_id);

CREATE TABLE blog_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    note_id UUID REFERENCES notes(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_blog_comments_blog_id ON blog_comments(blog_id);
CREATE INDEX idx_blog_comments_note_id ON blog_comments(note_id);
CREATE INDEX idx_blog_comments_user_id ON blog_comments(user_id);

-- ============================================
-- Theme Tables
-- ============================================

CREATE TABLE themes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    theme_url VARCHAR(500),
    css_url VARCHAR(500),
    preview_img VARCHAR(500),
    description TEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

-- ============================================
-- File Tables
-- ============================================

CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    size BIGINT NOT NULL,
    path VARCHAR(500),
    mime_type VARCHAR(100),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    album_id UUID REFERENCES albums(id) ON DELETE SET NULL,
    is_default_album BOOLEAN DEFAULT FALSE,
    from_file_id UUID REFERENCES files(id) ON DELETE SET NULL
);

CREATE INDEX idx_files_user_id ON files(user_id);

CREATE TABLE albums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255),
    seq INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_albums_user_id ON albums(user_id);
CREATE INDEX idx_albums_is_deleted ON albums(is_deleted);

-- ============================================
-- Note Images Table
-- ============================================

CREATE TABLE note_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    note_id UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    path VARCHAR(500),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_note_images_note_id ON note_images(note_id);
CREATE INDEX idx_note_images_user_id ON note_images(user_id);

-- ============================================
-- Config Tables
-- ============================================

CREATE TABLE configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value_str TEXT,
    value_arr JSONB,
    value_map JSONB,
    value_arr_map JSONB,
    is_arr BOOLEAN DEFAULT FALSE,
    is_map BOOLEAN DEFAULT FALSE,
    is_arr_map BOOLEAN DEFAULT FALSE,
    updated_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, key)
);

CREATE INDEX idx_configs_user_id ON configs(user_id);
CREATE INDEX idx_configs_key ON configs(key);

-- ============================================
-- Session Tables
-- ============================================

CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    data JSONB,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- ============================================
-- Token Tables
-- ============================================

CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tokens_token ON tokens(token);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);

-- ============================================
-- Other Tables
-- ============================================

CREATE TABLE email_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    to_email VARCHAR(255),
    subject VARCHAR(500),
    content TEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_logs_to_email ON email_logs(to_email);

CREATE TABLE suggestions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_suggestions_user_id ON suggestions(user_id);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    target_id UUID,
    target_type VARCHAR(50),
    reason TEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_user_id ON reports(user_id);
CREATE INDEX idx_reports_target ON reports(target_id, target_type);
