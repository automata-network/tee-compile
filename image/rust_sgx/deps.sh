cd /root && \
curl 'https://static.rust-lang.org/rustup/dist/x86_64-unknown-linux-gnu/rustup-init' --output /root/rustup-init && \
chmod +x /root/rustup-init && \
echo '1' | /root/rustup-init --default-toolchain ${rust_toolchain} && \
echo 'source /root/.cargo/env' >> /root/.bashrc && \
#/root/.cargo/bin/rustup component add rust-src rls rust-analysis clippy rustfmt && \
#/root/.cargo/bin/cargo install xargo && \
rm /root/rustup-init && rm -rf /root/.cargo/registry && rm -rf /root/.cargo/git 


cd /root && \
git clone --recursive https://github.com/intel/linux-sgx && \
cd linux-sgx && \
git checkout sgx_2.15.1 && \
./download_prebuilt.sh && \
make -j "$(nproc)" sdk_install_pkg && \
echo -e 'no\n/opt' | ./linux/installer/bin/sgx_linux_x64_sdk_2.15.101.1.bin && \
echo 'source /opt/sgxsdk/environment' >> /root/.bashrc && \
cd /root && \
rm -rf /root/linux-sgx